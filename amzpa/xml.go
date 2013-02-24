// Copyright 2012 Marco Dinacci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package amzpa provides functionality for using the
// Amazon Product Advertising service.

package amzpa

import (
	"encoding/xml"
)

type Image struct {
	XMLName xml.Name `xml:"MediumImage"`
	URL     string
	Height  uint16
	Width   uint16
}

type Item struct {
	XMLName       xml.Name `xml:"Item"`
	ASIN          string
	URL           string
	DetailPageURL string
	Author        string `xml:"ItemAttributes>Author"`
	Price         string `xml:"ItemAttributes>ListPrice>FormattedPrice"`
	PriceRaw      string `xml:"ItemAttributes>ListPrice>Amount"`
	MediumImage   Image
}

type ItemCollection struct {
	XMLName xml.Name `xml:"Items"`
	Items   []Item   `xml:"Item"`
	IsValid bool     `xml:"Request>IsValid"`
}

type ItemLookupResponse struct {
	XMLName        xml.Name `xml:"ItemLookupResponse"`
	ItemCollection ItemCollection
}

func unmarshal(contents []byte) (ItemLookupResponse, error) {
	itemLookupResponse := ItemLookupResponse{}
	err := xml.Unmarshal(contents, &itemLookupResponse)

	if err != nil {
		return ItemLookupResponse{}, err
	}

	return itemLookupResponse, err
}
