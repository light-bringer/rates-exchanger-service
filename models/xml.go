package models

import (
	"encoding/xml"
)

// Define the structs to match the updated XML structure

type Envelope struct {
	XMLName xml.Name       `xml:"Envelope"`
	Subject string         `xml:"subject"`
	Sender  Sender         `xml:"Sender"`
	Cube    StructuredCube `xml:"Cube"`
}

type Sender struct {
	Name string `xml:"name"`
}

// StructuredCube is the parent cube element containing date-specific cubes.
type StructuredCube struct {
	Cubes []DateCube `xml:"Cube"`
}

// DateCube represents a cube with a specific date, containing multiple currency rate entries.
type DateCube struct {
	Time    string      `xml:"time,attr"`
	Entries []RateEntry `xml:"Cube"`
}

type RateEntry struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
