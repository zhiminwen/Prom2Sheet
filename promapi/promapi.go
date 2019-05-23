package promapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/tealeg/xlsx"
)

func NewClient(promUrl, caFile, certFile, keyFile string) v1.API {
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Could not read CA file:%v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client cert and key file:%v", err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		},
	}
	client, err := api.NewClient(api.Config{
		Address:      promUrl,
		RoundTripper: transport,
	})
	if err != nil {
		log.Fatalf("Failed to create prometheus client: %v", err)
	}

	return v1.NewAPI(client)
}

func SaveSheet(sheet *xlsx.Sheet, promApi v1.API, sheetConfig Sheet) error {
	value, err := promApi.Query(context.Background(), sheetConfig.Query, time.Now())
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()
	for _, col := range sheetConfig.Columns {
		cell = row.AddCell()
		cell.SetString(col.Name)
	}

	for _, v := range value.(model.Vector) {
		row = sheet.AddRow()
		for _, col := range sheetConfig.Columns {
			cell = row.AddCell()
			switch col.Type {
			case "OS.Environment":
				cell.SetString(os.Getenv(col.Value))
			case "Prometheus.Timestamp":
				cell.SetDateTime(v.Timestamp.Time())
			case "Prometheus.Metric":
				cell.SetValue(string(v.Metric[model.LabelName(col.Value)]))
			case "Prometheus.Value":
				cell.SetFloat(float64(v.Value))
			}
		}
	}

	return nil
}
