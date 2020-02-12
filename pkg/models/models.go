package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Website struct {
	ID    int
	Title string
}

type Sub struct {
	ID        int
	Title     string
	WebsiteId int
}

type Pcap struct {
	ID    int
	Path  string
	SubId int
	Proxy string
}

type Feature struct {
	ID             int
	PcapID         int
	FeatureSet     []byte
	FeatureSetJson FeatureSetJson
}

type FeatureSetJson struct {
	Cumul [50]int `json:"int"`
}
