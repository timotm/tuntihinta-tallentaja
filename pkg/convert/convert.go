package convert

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

type TimeInterval struct {
	Start string `xml:"start"`
	End   string `xml:"end"`
}

type Point struct {
	Position int     `xml:"position"`
	Price    float64 `xml:"price.amount"`
}

type Period struct {
	TimeInterval TimeInterval `xml:"timeInterval"`
	Resolution   string       `xml:"resolution"`
	Point        []Point      `xml:"Point"`
}

type TimeSeries struct {
	MRId             string `xml:"mRID"`
	BusinessType     string `xml:"businessType"`
	InMRId           string `xml:"in_Domain.mRID"`
	OutMRId          string `xml:"out_Domain.mRID"`
	Currency         string `xml:"currency_Unit.name"`
	PriceMeasureUnit string `xml:"price_Measure_Unit.name"`
	CurveType        string `xml:"curveType"`
	Period           Period `xml:"Period"`
}

type MarketDocument struct {
	MRId            string       `xml:"mRID"`
	RevisionNumber  string       `xml:"revisionNumber"`
	Type            string       `xml:"type"`
	SenderMRId      string       `xml:"sender_MarketParticipant.mRID"`
	SenderType      string       `xml:"sender_MarketParticipant.marketRole.type"`
	ReceiverMRId    string       `xml:"receiver_MarketParticipant.mRID"`
	ReceiverType    string       `xml:"receiver_MarketParticipant.marketRole.type"`
	CreatedDateTime string       `xml:"createdDateTime"`
	TimeInterval    TimeInterval `xml:"period.timeInterval"`
	TimeSeries      []TimeSeries `xml:"TimeSeries"`
}

type Reason struct {
	Code string `xml:"code"`
	Text string `xml:"text"`
}
type Acknowledgement_MarketDocument struct {
	Reason Reason `xml:"Reason"`
}

type HourPrice struct {
	StartTime time.Time
	Price     float64
}

func (hp *HourPrice) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf(`{"startTime":"%s","price":%.3f}`,
		hp.StartTime.Format("2006-01-02T15:04:05Z"),
		hp.Price)
	return []byte(json), nil
}

type DayPrices struct {
	Date       string      `json:"date"`
	HourPrices []HourPrice `json:"hourPrices"`
}

func parseXMLTime(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04Z", input)
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("Unable to parse time %s: %s", input, err))
	}
	return t, nil
}

func ParseXml(xmlBytes []byte) ([]DayPrices, error) {
	var marketDocument MarketDocument
	if err := xml.Unmarshal(xmlBytes, &marketDocument); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse XML: %s", err))
	}

	if len(marketDocument.TimeSeries) == 0 {
		var ack Acknowledgement_MarketDocument
		if err := xml.Unmarshal(xmlBytes, &ack); err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to parse XML (ack): %s", err))
		}
		return nil, errors.New(fmt.Sprintf("No time series found: %s / %s", ack.Reason.Code, ack.Reason.Text))
	}

	var dayPrices []DayPrices

	for i := range marketDocument.TimeSeries {
		start, err := parseXMLTime(marketDocument.TimeSeries[i].Period.TimeInterval.Start)
		if err != nil {
			return nil, err
		}

		var minutes int

		if marketDocument.TimeSeries[i].Period.Resolution == "PT60M" {
			minutes = 60
		} else if marketDocument.TimeSeries[i].Period.Resolution == "PT15M" {
			minutes = 15
		} else {
			return nil, errors.New(fmt.Sprintf("Unsupported resolution %s", marketDocument.TimeSeries[i].Period.Resolution))
		}

		if marketDocument.TimeSeries[i].Currency != "EUR" {
			return nil, errors.New(fmt.Sprintf("Unsupported currency %s", marketDocument.TimeSeries[i].Currency))
		}

		if marketDocument.TimeSeries[i].PriceMeasureUnit != "MWH" {
			return nil, errors.New(fmt.Sprintf("Unsupported price measure unit %s", marketDocument.TimeSeries[i].PriceMeasureUnit))
		}

		for j := range marketDocument.TimeSeries[i].Period.Point {
			if j == 0 {
				dayPrices = append(dayPrices,
					DayPrices{Date: start.Round(24 * time.Hour).Format("2006-01-02"),
						HourPrices: []HourPrice{}})
			}
			currentTime := start.Add(time.Minute * time.Duration(minutes*j))
			dayPrices[len(dayPrices)-1].HourPrices = append(dayPrices[len(dayPrices)-1].HourPrices,
				HourPrice{StartTime: currentTime,
					Price: marketDocument.TimeSeries[i].Period.Point[j].Price / 10.0}) // e/MWh -> c/KWh

		}
	}

	return dayPrices, nil
}
