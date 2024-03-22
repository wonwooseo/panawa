package kamis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/wonwooseo/panawa/pkg/db/model"
)

type DataClient struct {
	logger zerolog.Logger
	client *http.Client

	apiURL string
}

func NewDataClient(baseLogger zerolog.Logger, apiURL string) *DataClient {
	return &DataClient{
		logger: baseLogger.With().Str("caller", "data/kamis").Logger(),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		apiURL: apiURL,
	}
}

func (c *DataClient) GetDatePrices(ctx context.Context, date time.Time, itemCode string) (*model.Price, map[string][]*model.Price, error) {
	productCodes, ok := kamisProductCodeMap[itemCode]
	if !ok {
		return nil, nil, fmt.Errorf("unknown item code: %s", itemCode)
	}

	// 소매가격 정보 취득
	reqURL, err := url.Parse(c.apiURL)
	if err != nil {
		return nil, nil, err
	}
	query := reqURL.Query()
	query.Add("p_productclscode", "01")
	query.Add("p_startday", date.Format("2006-01-02"))
	query.Add("p_endday", date.Format("2006-01-02"))
	query.Add("p_itemcategorycode", productCodes.CategoryCode)
	query.Add("p_itemcode", productCodes.ItemCode)
	query.Add("p_kindcode", productCodes.KindCode)
	query.Add("p_productrankcode", productCodes.RankCode)
	query.Add("p_convert_kg_yn", "N")
	query.Add("p_cert_key", "111") // api 이용신청 후 변경
	query.Add("p_cert_id", "222")  // api 이용신청 후 변경
	query.Add("p_returntype", "json")
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	c.logger.Info().Str("request_url", reqURL.String()).Msg("sending price data request")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	responseTime := time.Now().UTC()
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	parsed := apiResponse{}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, nil, err
	}

	if parsed.Data.ErrorCode != statusSuccess {
		switch parsed.Data.ErrorCode {
		case statusUnauthenticated:
			return nil, nil, fmt.Errorf("unauthenticated price data request")
		case statusWrongParameters:
			return nil, nil, fmt.Errorf("invalid price data request parameters")
		default:
			c.logger.Error().Any("response", parsed).Msg("unknown error code")
			return nil, nil, fmt.Errorf("unknown error code from price data response")
		}
	}

	type tempPrice struct {
		Low   int
		Sum   int
		High  int
		Count int
	}
	tempDatePrice := tempPrice{Low: math.MaxInt, Sum: 0, High: 0, Count: 0}
	tempRegionalMarketPrices := map[string]map[string]tempPrice{}
	for _, price := range parsed.Data.Item {
		regionCode, ok := regionCodeMap[price.CountyName.String()]
		if !ok {
			continue
		}
		marketCode := "00" // 전통시장
		if isSuperMarket(price.MarketName.String()) {
			marketCode = "01" // 대형유통
		}
		priceInt, err := strconv.Atoi(strings.ReplaceAll(price.Price, ",", ""))
		if err != nil {
			c.logger.Warn().Any("data", price).Err(err).Msg("failed to convert price from response data to int")
			continue
		}
		c.logger.Info().Str("region", regionCode).Str("market", marketCode).Int("price", priceInt).Msg("price data")

		// for date price
		if priceInt < tempDatePrice.Low {
			tempDatePrice.Low = priceInt
		}
		if priceInt > tempDatePrice.High {
			tempDatePrice.High = priceInt
		}
		tempDatePrice.Sum += priceInt
		tempDatePrice.Count += 1

		// for regional market prices
		_, ok = tempRegionalMarketPrices[regionCode]
		if !ok {
			tempRegionalMarketPrices[regionCode] = map[string]tempPrice{}
		}
		rmp, ok := tempRegionalMarketPrices[regionCode][marketCode]
		if !ok {
			tempRegionalMarketPrices[regionCode][marketCode] = tempPrice{
				Low:   priceInt,
				Sum:   priceInt,
				High:  priceInt,
				Count: 1,
			}
		} else {
			if priceInt < rmp.Low {
				rmp.Low = priceInt
			}
			if priceInt > rmp.High {
				rmp.High = priceInt
			}
			rmp.Sum += priceInt
			rmp.Count += 1
			tempRegionalMarketPrices[regionCode][marketCode] = rmp
		}
	}

	datePrice := &model.Price{
		ItemCode:   itemCode,
		Low:        tempDatePrice.Low,
		Average:    int(tempDatePrice.Sum / tempDatePrice.Count),
		High:       tempDatePrice.High,
		UpdateTime: responseTime,
	}
	regionalMarketPrices := map[string][]*model.Price{}
	for regionCode, marketCodeTempPriceMap := range tempRegionalMarketPrices {
		for marketCode, tempPrice := range marketCodeTempPriceMap {
			regionalMarketPrices[regionCode] = append(regionalMarketPrices[regionCode], &model.Price{
				ItemCode:   itemCode,
				Low:        tempPrice.Low,
				Average:    int(tempPrice.Sum / tempPrice.Count),
				High:       tempPrice.High,
				RegionCode: &regionCode, // safe to take pointer of loop var after go 1.22
				MarketCode: &marketCode, // safe to take pointer of loop var after go 1.22
				UpdateTime: responseTime,
			})
		}
	}
	return datePrice, regionalMarketPrices, nil
}
