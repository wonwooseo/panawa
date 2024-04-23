package kamis

import "encoding/json"

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
"data": [ // 데이터가 없는 경우(주말, 아직 업데이트 안됨) 배열이 내려온다
	"001"
]
*/

type kamisDataOrArray []interface{}

type kamisData struct {
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
	Valid bool `json:"-"` // internally used to check validity
}

func (a *kamisDataOrArray) UnmarshalJSON(b []byte) error {
	switch b[0] {
	case '[':
		var strArr []string
		if err := json.Unmarshal(b, &strArr); err != nil {
			return err
		}
		t := make([]interface{}, len(strArr))
		for i, s := range strArr {
			t[i] = s
		}
		*a = t
	default:
		var data kamisData
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
		t := make([]interface{}, 1)
		t[0] = data
		*a = t
	}
	return nil
}

func (a kamisDataOrArray) KamisData() kamisData {
	t := []interface{}(a)
	if len(t) > 0 {
		data, ok := t[0].(kamisData)
		if !ok {
			return kamisData{}
		}
		data.Valid = true
		return data
	}
	return kamisData{}
}

type apiResponse struct {
	Data kamisDataOrArray `json:"data"`
}
