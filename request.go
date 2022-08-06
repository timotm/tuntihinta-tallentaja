package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getXmlPriceDataForDateRange(start string, end string, securityToken string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://transparency.entsoe.eu/api", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("securityToken", securityToken)
	q.Add("documentType", "A44")
	q.Add("in_domain", "10YFI-1--------U")
	q.Add("out_domain", "10YFI-1--------U")
	q.Add("TimeInterval", fmt.Sprintf("%s/%s", start, end))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Error %d: %s", resp.StatusCode, responseBody)
	}

	return responseBody
}
