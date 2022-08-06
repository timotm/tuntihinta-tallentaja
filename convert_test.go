package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
	<Publication_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-3:publicationdocument:7:0">
		<mRID>4d71ec7ac22542abafb6fb1bb912fa34</mRID>
		<revisionNumber>1</revisionNumber>
		<type>A44</type>
		<sender_MarketParticipant.mRID codingScheme="A01">10X1001A1001A450</sender_MarketParticipant.mRID>
		<sender_MarketParticipant.marketRole.type>A32</sender_MarketParticipant.marketRole.type>
		<receiver_MarketParticipant.mRID codingScheme="A01">10X1001A1001A450</receiver_MarketParticipant.mRID>
		<receiver_MarketParticipant.marketRole.type>A33</receiver_MarketParticipant.marketRole.type>
		<createdDateTime>2022-08-06T12:38:33Z</createdDateTime>
		<period.timeInterval>
			<start>2021-12-31T23:00Z</start>
			<end>2022-08-07T22:00Z</end>
		</period.timeInterval>
		<TimeSeries>
			<mRID>1</mRID>
			<businessType>A62</businessType>
			<in_Domain.mRID codingScheme="A01">10YFI-1--------U</in_Domain.mRID>
			<out_Domain.mRID codingScheme="A01">10YFI-1--------U</out_Domain.mRID>
			<currency_Unit.name>EUR</currency_Unit.name>
			<price_Measure_Unit.name>MWH</price_Measure_Unit.name>
			<curveType>A01</curveType>
			<Period>
				<timeInterval>
					<start>2021-12-31T23:00Z</start>
					<end>2022-01-01T23:00Z</end>
				</timeInterval>
				<resolution>PT60M</resolution>
				<Point>
					<position>1</position>
					<price.amount>46.60</price.amount>
				</Point>
				<Point>
					<position>2</position>
					<price.amount>26.60</price.amount>
				</Point>
			</Period>
		</TimeSeries>
		<TimeSeries>
			<mRID>2</mRID>
			<businessType>A62</businessType>
			<in_Domain.mRID codingScheme="A01">10YFI-1--------U</in_Domain.mRID>
			<out_Domain.mRID codingScheme="A01">10YFI-1--------U</out_Domain.mRID>
			<currency_Unit.name>EUR</currency_Unit.name>
			<price_Measure_Unit.name>MWH</price_Measure_Unit.name>
			<curveType>A01</curveType>
			<Period>
				<timeInterval>
					<start>2022-01-01T23:00Z</start>
					<end>2022-01-02T23:00Z</end>
				</timeInterval>
				<resolution>PT60M</resolution>
				<Point>
					<position>1</position>
					<price.amount>57.08</price.amount>
				</Point>
				<Point>
					<position>2</position>
					<price.amount>55.09</price.amount>
				</Point>
				<Point>
					<position>3</position>
					<price.amount>52.82</price.amount>
				</Point>
			</Period>
		</TimeSeries>
	</Publication_MarketDocument>`

	prices, err := parseXml([]byte(input))
	assert.Nil(t, err)

	pricesJson, _ := json.Marshal(prices)

	expectedPretty := `[
		{
			"date": "2022-01-01",
			"hourPrices": [
				{
					"startTime": "2021-12-31T23:00:00Z",
					"price": 4.660
				},
				{
					"startTime": "2022-01-01T00:00:00Z",
					"price": 2.660
				}
			]
		},
		{
			"date": "2022-01-02",
			"hourPrices": [
				{
					"startTime": "2022-01-01T23:00:00Z",
					"price": 5.708
				},
				{
					"startTime": "2022-01-02T00:00:00Z",
					"price": 5.509
				},
				{
					"startTime": "2022-01-02T01:00:00Z",
					"price": 5.282
				}
			]
		}
	]`
	expected := &bytes.Buffer{}
	if err := json.Compact(expected, []byte(expectedPretty)); err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected.String(), string(pricesJson))
}

func TestNoData(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
	<Acknowledgement_MarketDocument
		xmlns="urn:iec62325.351:tc57wg16:451-1:acknowledgementdocument:7:0">
		<mRID>3273dac7-459c-4</mRID>
		<createdDateTime>2022-08-07T10:01:42Z</createdDateTime>
		<sender_MarketParticipant.mRID codingScheme="A01">10X1001A1001A450</sender_MarketParticipant.mRID>
		<sender_MarketParticipant.marketRole.type>A32</sender_MarketParticipant.marketRole.type>
		<receiver_MarketParticipant.mRID codingScheme="A01">10X1001A1001A450</receiver_MarketParticipant.mRID>
		<receiver_MarketParticipant.marketRole.type>A39</receiver_MarketParticipant.marketRole.type>
		<received_MarketDocument.createdDateTime>2022-08-07T10:01:42Z</received_MarketDocument.createdDateTime>
		<Reason>
			<code>999</code>
			<text>No matching data found for Data item Day-ahead Prices [12.1.D] (10YFI-1--------U, 10YFI-1--------U) and interval 2022-08-08T00:00:00.000Z/2022-08-08T00:00:00.000Z.</text>
		</Reason>
	</Acknowledgement_MarketDocument>`

	_, err := parseXml([]byte(input))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "No matching data found for Data item Day-ahead")
}
