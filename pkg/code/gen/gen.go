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

	lang := viper.GetString("language")
	item := viper.GetStringMapString("item")
	region := viper.GetStringMapString("region")
	market := viper.GetStringMapString("market")

	tmpl := template.Must(template.ParseFiles("resolver.tmpl"))

	if err := os.MkdirAll(lang, os.ModePerm); err != nil {
		panic(err)
	}

	itemGo, err := os.Create(fmt.Sprintf("%s/item.go", lang))
	if err != nil {
		panic(err)
	}
	defer itemGo.Close()
	var itemBuf bytes.Buffer
	if err := tmpl.Execute(&itemBuf, templateArg{
		Pkg:           lang,
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

	regionGo, err := os.Create(fmt.Sprintf("%s/region.go", lang))
	if err != nil {
		panic(err)
	}
	defer regionGo.Close()
	var regionBuf bytes.Buffer
	if err := tmpl.Execute(&regionBuf, templateArg{
		Pkg:           lang,
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

	marketGo, err := os.Create(fmt.Sprintf("%s/market.go", lang))
	if err != nil {
		panic(err)
	}
	defer marketGo.Close()
	var marketBuf bytes.Buffer
	if err := tmpl.Execute(&marketBuf, templateArg{
		Pkg:           lang,
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
