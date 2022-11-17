Fetches Finnish day-ahead electricity prices from [ENTSO-E transparency platform](https://transparency.entsoe.eu/transmission-domain/r2/dayAheadPrices/show?name=&defaultValue=true&viewType=GRAPH&areaType=BZN&atch=false&dateTime.dateTime=08.08.2022+00:00|EET|DAY&biddingZone.values=CTY|10YFI-1--------U!BZN|10YFI-1--------U&resolution.values=PT15M&resolution.values=PT30M&resolution.values=PT60M&dateTime.timezone=EET_EEST&dateTime.timezone_input=EET+(UTC+2)+/+EEST+(UTC+3)) with an API key and stores them in a simpler JSON format in S3.

Build and run:
```
go build && TH_AWS_REGION=eu-north-1 TH_AWS_BUCKET_NAME=buket TH_AWS_ACCESS_KEY_ID=yes TH_AWS_SECRET_ACCESS_KEY=very_secret TH_SECURITY_TOKEN=also_secret ./tuntihinta-tallentaja --start-date 2022-08-01 --end-date 2022-08-07
```

Sample JSON file:
```JSON
{
  "date": "2022-08-08",
  "hourPrices": [
    { "startTime": "2022-08-07T22:00:00Z", "price": 10.002 },
    { "startTime": "2022-08-07T23:00:00Z", "price": 8.001 },
    { "startTime": "2022-08-08T00:00:00Z", "price": 8.004 },
    { "startTime": "2022-08-08T01:00:00Z", "price": 11.87 },
    { "startTime": "2022-08-08T02:00:00Z", "price": 15.353 },
    { "startTime": "2022-08-08T03:00:00Z", "price": 31.038 },
    { "startTime": "2022-08-08T04:00:00Z", "price": 46.289 },
    { "startTime": "2022-08-08T05:00:00Z", "price": 79.832 },
    { "startTime": "2022-08-08T06:00:00Z", "price": 86.114 },
    { "startTime": "2022-08-08T07:00:00Z", "price": 57.409 },
    { "startTime": "2022-08-08T08:00:00Z", "price": 44.396 },
    { "startTime": "2022-08-08T09:00:00Z", "price": 75.003 },
    { "startTime": "2022-08-08T10:00:00Z", "price": 79.219 },
    { "startTime": "2022-08-08T11:00:00Z", "price": 79.798 },
    { "startTime": "2022-08-08T12:00:00Z", "price": 44.793 },
    { "startTime": "2022-08-08T13:00:00Z", "price": 47.896 },
    { "startTime": "2022-08-08T14:00:00Z", "price": 46.693 },
    { "startTime": "2022-08-08T15:00:00Z", "price": 47.964 },
    { "startTime": "2022-08-08T16:00:00Z", "price": 86.111 },
    { "startTime": "2022-08-08T17:00:00Z", "price": 79.799 },
    { "startTime": "2022-08-08T18:00:00Z", "price": 49.522 },
    { "startTime": "2022-08-08T19:00:00Z", "price": 45.5 },
    { "startTime": "2022-08-08T20:00:00Z", "price": 42.915 },
    { "startTime": "2022-08-08T21:00:00Z", "price": 31.171 }
  ]
}
```
