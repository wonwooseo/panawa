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

func isSuperMarket(marketName string) bool {
	return strings.Contains(marketName, "유통")
}
