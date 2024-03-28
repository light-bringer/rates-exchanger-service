package sync

import "encoding/xml"

type XMLExchangeRate struct {
	XMLName xml.Name   `xml:"Envelope"`
	Cube    CubeStruct `xml:"http://www.ecb.int/vocabulary/2002-08-01/eurofxref Cube"`
}

type CubeStruct struct {
	Time        string `xml:"time,attr"`
	CubeContent []Cube `xml:"Cube"`
}

type Cube struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
	Time     string  `xml:"time,attr"`
}

type (
	ExchangeRate  Cube
	ExchangeRates []ExchangeRate
)
