package kamis

import "strings"

type productCode struct {
	CategoryCode string
	ItemCode     string
	KindCode     string
	RankCode     string
}

// MUST match codes from `pkg/code`
var kamisProductCodeMap = map[string]productCode{
	"0000": {
		CategoryCode: "200",
		ItemCode:     "246",
		KindCode:     "00",
		RankCode:     "04",
	},
}

const (
	statusSuccess         = "000"
	statusUnauthenticated = "900"
	statusWrongParameters = "200"
)

// MUST match codes from `pkg/code`
var regionCodeMap = map[string]string{
	"서울": "0000",
	"부산": "0001",
	"대구": "0002",
	"인천": "0003",
	"광주": "0004",
	"대전": "0005",
	"울산": "0006",
	"세종": "0007",
	"수원": "0008",
	"성남": "0009",
	"고양": "0010",
	"용인": "0011",
	"춘천": "0012",
	"강릉": "0013",
	"청주": "0014",
	"천안": "0015",
	"전주": "0016",
	"순천": "0017",
	"포항": "0018",
	"안동": "0019",
	"창원": "0020",
	"김해": "0021",
	"제주": "0022",
}

func isSuperMarket(marketName string) bool {
	return strings.Contains(marketName, "유통")
}
