package glue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/timotm/tuntihinta-tallentaja/pkg/convert"
	"github.com/timotm/tuntihinta-tallentaja/pkg/fetcher"
	"github.com/timotm/tuntihinta-tallentaja/pkg/s3"
)

func FetchAndUpload(startTime, endTime fetcher.EntsoeTime,
	token, region, bucket, accessKey, secretAccessKey string) []string {
	responseXml := fetcher.GetXmlPriceDataForDateRange(startTime.String(), endTime.String(), token)

	prices, err := convert.ParseXml(responseXml)
	if err != nil {
		log.Fatalf("Unable to parse XML: %s", err)
	}

	files := make([]string, 0, len(prices))

	for _, dayPrices := range prices {
		fileName := fmt.Sprintf("%s.json", dayPrices.Date)
		contents, err := json.Marshal(dayPrices)
		if err != nil {
			log.Fatalf("Unable to marshal JSON: %s", err)
		}
		files = append(files, fileName)
		s3.PutFileToS3(contents, region, accessKey, secretAccessKey, bucket, fileName)
	}

	return files
}
