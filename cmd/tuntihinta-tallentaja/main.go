package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/timotm/tuntihinta-tallentaja/pkg/fetcher"
	"github.com/timotm/tuntihinta-tallentaja/pkg/glue"
)

var startTime fetcher.EntsoeTime
var endTime fetcher.EntsoeTime

func init() {
	flag.Var(&startTime, "start-date", "First date to fetch data for, e.g. 2021-01-01. Default tomorrow.")
	flag.Var(&endTime, "end-date", "Last date to fetch data for, default tomorrow.")
}

func main() {
	flag.Parse()

	requiredEnvVariables := []string{"TH_SECURITY_TOKEN", "TH_AWS_REGION", "TH_AWS_BUCKET_NAME", "TH_AWS_ACCESS_KEY_ID", "TH_AWS_SECRET_ACCESS_KEY"}
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
	files := glue.FetchAndUpload(startTime,
		endTime,
		os.Getenv("TH_SECURITY_TOKEN"),
		os.Getenv("TH_AWS_REGION"),
		os.Getenv("TH_AWS_BUCKET_NAME"),
		os.Getenv("TH_AWS_ACCESS_KEY_ID"),
		os.Getenv("TH_AWS_SECRET_ACCESS_KEY"))

	fmt.Printf("Uploaded files: %s\n", strings.Join(files, ", "))
}
