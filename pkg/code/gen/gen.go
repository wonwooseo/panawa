//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

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

	country := viper.GetString("country")
	lang := viper.GetString("language")
	item := viper.GetStringMapString("item")
	region := viper.GetStringMapString("region")
	market := viper.GetStringMapString("market")

	tmpl := template.Must(template.ParseFiles("resolver.tmpl"))

	if err := os.MkdirAll(fmt.Sprintf("%s/%s", country, lang), os.ModePerm); err != nil {
		panic(err)
	}

	itemGo, err := os.Create(fmt.Sprintf("%s/%s/item.go", country, lang))
	if err != nil {
		panic(err)
	}
	defer itemGo.Close()
	var itemBuf bytes.Buffer
	if err := tmpl.Execute(&itemBuf, templateArg{
		Pkg:           fmt.Sprintf("%s%s", country, lang),
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

	regionGo, err := os.Create(fmt.Sprintf("%s/%s/region.go", country, lang))
	if err != nil {
		panic(err)
	}
	defer regionGo.Close()
	var regionBuf bytes.Buffer
	if err := tmpl.Execute(&regionBuf, templateArg{
		Pkg:           fmt.Sprintf("%s%s", country, lang),
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

	marketGo, err := os.Create(fmt.Sprintf("%s/%s/market.go", country, lang))
	if err != nil {
		panic(err)
	}
	defer marketGo.Close()
	var marketBuf bytes.Buffer
	if err := tmpl.Execute(&marketBuf, templateArg{
		Pkg:           fmt.Sprintf("%s%s", country, lang),
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
