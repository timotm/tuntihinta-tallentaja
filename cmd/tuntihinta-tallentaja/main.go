package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/timotm/tuntihinta-tallentaja/pkg/convert"
	"github.com/timotm/tuntihinta-tallentaja/pkg/fetcher"
	"github.com/timotm/tuntihinta-tallentaja/pkg/s3"
)

var startTime fetcher.EntsoeTime
var endTime fetcher.EntsoeTime

func init() {
	flag.Var(&startTime, "start-date", "First date to fetch data for, e.g. 2021-01-01. Default tomorrow.")
	flag.Var(&endTime, "end-date", "Last date to fetch data for, default tomorrow.")
}

func main() {
	flag.Parse()

	requiredEnvVariables := []string{"SECURITY_TOKEN", "AWS_REGION", "AWS_BUCKET_NAME", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}
	var unset []string
	for _, envVar := range requiredEnvVariables {
		if os.Getenv(envVar) == "" {
			unset = append(unset, envVar)
		}
	}
	if len(unset) > 0 {
		log.Fatalf("Missing environment variables: %s", strings.Join(unset, ", "))
	}

	if !startTime.IsSet() {
		t := time.Now().AddDate(0, 0, 1)
		startTime = fetcher.EntsoeTime(t)
	}

	if !endTime.IsSet() {
		t := time.Now().AddDate(0, 0, 2)
		endTime = fetcher.EntsoeTime(t)
	}

	fmt.Printf("Fetching data for %+v to %+v\n", startTime.String(), endTime.String())
	responseXml := fetcher.GetXmlPriceDataForDateRange(startTime.String(), endTime.String(), os.Getenv("SECURITY_TOKEN"))

	prices, err := convert.ParseXml(responseXml)
	if err != nil {
		log.Fatalf("Unable to parse XML: %s", err)
	}

	for _, dayPrices := range prices {
		fileName := fmt.Sprintf("%s.json", dayPrices.Date)
		contents, err := json.Marshal(dayPrices)
		if err != nil {
			log.Fatalf("Unable to marshal JSON: %s", err)
		}
		fmt.Printf("Writing %s: %s\n", fileName, contents)
		s3.PutFileToS3(contents, os.Getenv("AWS_REGION"), os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_BUCKET_NAME"), fileName)
	}
}
