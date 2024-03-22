package kamis

import "encoding/json"

type stringOrArray []string

func (a *stringOrArray) UnmarshalJSON(b []byte) error {
	switch b[0] {
	case '[':
		var strArr []string
		if err := json.Unmarshal(b, &strArr); err != nil {
			return err
		}
		*a = stringOrArray(strArr)
	default:
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		*a = stringOrArray{str}
	}
	return nil
}

func (a stringOrArray) String() string {
	if len([]string(a)) > 0 {
		return a[0]
	}
	return ""
}

/*
"data": {
    "error_code": "000",
    "item": [
        {
			"itemname": [], // 배열일 수도 있고
			"kindname": [],
			"countyname": "평균",
			"marketname": [],
			"yyyy": "2024",
			"regday": "03/20",
			"price": "2,878"
		},
		{
			"itemname": "파", // 문자열일 수도 있다
			"kindname": "대파(1kg)",
			"countyname": "서울",
			"marketname": "경동",
			"yyyy": "2024",
			"regday": "03/20",
			"price": "4,000"
		},
		...
	]
}
*/
type apiResponse struct {
	Data struct {
		ErrorCode string `json:"error_code"`
		Item      []struct {
			ItemName   stringOrArray `json:"itemname"`
			KindName   stringOrArray `json:"kindname"`
			CountyName stringOrArray `json:"countyname"`
			MarketName stringOrArray `json:"marketname"`
			Year       string        `json:"yyyy"`
			RegDay     string        `json:"regday"`
			Price      string        `json:"price"`
		} `json:"item"`
	} `json:"data"`
}
