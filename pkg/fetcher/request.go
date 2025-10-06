package fetcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type EntsoeTime time.Time

func (t *EntsoeTime) String() string {
	return time.Time(*t).Format("200601020000")
}

func (t *EntsoeTime) Set(value string) error {
	tmp, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*t = EntsoeTime(tmp)
	return nil
}

func (t *EntsoeTime) IsSet() bool {
	return !time.Time(*t).IsZero()
}

func GetXmlPriceDataForDateRange(start string, end string, securityToken string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://web-api.tp.entsoe.eu/api", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("securityToken", securityToken)
	q.Add("documentType", "A44")
	q.Add("in_Domain", "10YFI-1--------U")
	q.Add("out_Domain", "10YFI-1--------U")
	q.Add("periodStart", start)
	q.Add("periodEnd", end)
	req.URL.RawQuery = q.Encode()

	fmt.Printf("Requesting URL: %s\n", req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Error %d: %s", resp.StatusCode, responseBody)
	}

	return responseBody
}
