//go:build ignore

package main

import (
	"bytes"
	"go/format"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

const pkgName = "code"

type templateArg struct {
	Pkg           string
	Type          string
	CodeLocaleMap map[string]string
}

func main() {
	cfgPath, ok := os.LookupEnv("GEN_CONFIG")
	if !ok {
		panic("gen config path not provided")
	}

	viper.SetConfigFile(cfgPath)
	viper.ReadInConfig()

	item := viper.GetStringMapString("item")
	region := viper.GetStringMapString("region")
	market := viper.GetStringMapString("market")

	tmpl := template.Must(template.ParseFiles("resolver.tmpl"))

	itemGo, err := os.Create("./item.go")
	if err != nil {
		panic(err)
	}
	defer itemGo.Close()
	var itemBuf bytes.Buffer
	if err := tmpl.Execute(&itemBuf, templateArg{
		Pkg:           pkgName,
		Type:          "Item",
		CodeLocaleMap: item,
	}); err != nil {
		panic(err)
	}
	itemFmt, err := format.Source(itemBuf.Bytes())
	if err != nil {
		panic(err)
	}
	itemGo.Write(itemFmt)

	regionGo, err := os.Create("./region.go")
	if err != nil {
		panic(err)
	}
	defer regionGo.Close()
	var regionBuf bytes.Buffer
	if err := tmpl.Execute(&regionBuf, templateArg{
		Pkg:           pkgName,
		Type:          "Region",
		CodeLocaleMap: region,
	}); err != nil {
		panic(err)
	}
	regionFmt, err := format.Source(regionBuf.Bytes())
	if err != nil {
		panic(err)
	}
	regionGo.Write(regionFmt)

	marketGo, err := os.Create("./market.go")
	if err != nil {
		panic(err)
	}
	defer marketGo.Close()
	var marketBuf bytes.Buffer
	if err := tmpl.Execute(&marketBuf, templateArg{
		Pkg:           pkgName,
		Type:          "Market",
		CodeLocaleMap: market,
	}); err != nil {
		panic(err)
	}
	marketFmt, err := format.Source(marketBuf.Bytes())
	if err != nil {
		panic(err)
	}
	marketGo.Write(marketFmt)
}
