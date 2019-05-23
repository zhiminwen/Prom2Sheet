// +build mage

package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/tealeg/xlsx"

	"./promapi"
)

type Configuration struct {
	PrometheusAddress string `split_words:"true" default:"https://monitoring-prometheus:9090"`
	CaFile            string `split_words:"true" default:"certs/myca.pem"`
	CertFile          string `split_words:"true" default:"certs/prom.pem"`
	KeyFile           string `split_words:"true" default:"certs/prom-key.pem"`
	SheetYaml         string `split_words:"true" default:"config.yaml"`
}

var config Configuration
var promApi v1.API
var sheetConfig *promapi.XLSConfig

func init() {
	os.Setenv("MAGEFILE_VERBOSE", "true")

	err := envconfig.Process("P", &config)
	if err != nil {
		log.Fatalf("Failed to load config file:%v", err)
	}
	log.Printf("config:%v", config)

	promApi = promapi.NewClient(config.PrometheusAddress, config.CaFile, config.CertFile, config.KeyFile)
	sheetConfig, err = promapi.ParseSheetYaml(config.SheetYaml)
	if err != nil {
		log.Fatalf("Failed to parse sheet yaml:%v", err)
	}
}

func Query() {
	xlsFile := xlsx.NewFile()

	for _, shtConf := range sheetConfig.Sheets {
		sheet, err := xlsFile.AddSheet(shtConf.Name)
		if err != nil {
			log.Fatalf("Failed to add sheet:%v", err)
		}
		err = promapi.SaveSheet(sheet, promApi, shtConf)
		if err != nil {
			log.Fatalf("Failed to save the sheet:%v", err)
		}
	}
	err := xlsFile.Save("./metering.xlsx")
	if err != nil {
		log.Fatalf("Failed to save xls file:%v", err)
	}
}
