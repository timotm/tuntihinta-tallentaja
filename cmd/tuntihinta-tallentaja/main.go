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

type entsoeTime time.Time

var startTime entsoeTime
var endTime entsoeTime

func (t *entsoeTime) String() string {
	return time.Time(*t).Format("2006-01-02T00:00Z")
}

func (t *entsoeTime) Set(value string) error {
	tmp, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*t = entsoeTime(tmp)
	return nil
}

func (t *entsoeTime) isSet() bool {
	return !time.Time(*t).IsZero()
}

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

	if !startTime.isSet() {
		t := time.Now().AddDate(0, 0, 1)
		startTime = entsoeTime(t)
	}

	if !endTime.isSet() {
		t := time.Now().AddDate(0, 0, 2)
		endTime = entsoeTime(t)
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
		s3.PutFileToS3(contents, os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET_NAME"), fileName)
	}
}
