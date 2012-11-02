// Copyright 2012 Marco Dinacci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package amzpa provides functionality for using the
// Amazon Product Advertising service.

package amzpa

import (
	"fmt"
	"sort"
	"time"
	"io/ioutil"
	"strings"
	"net/url"
	"net/http"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
)

var service_domains = map[string] string {
     "CA": "ecs.amazonaws.ca",
     "CN": "webservices.amazon.cn",
     "DE": "ecs.amazonaws.de",
     "ES": "webservices.amazon.es",
     "FR": "ecs.amazonaws.fr",
     "IT": "webservices.amazon.it",
     "JP": "ecs.amazonaws.jp",
     "UK": "ecs.amazonaws.co.uk",
     "US": "ecs.amazonaws.com",
}

type AmazonRequest struct {
	accessKeyID string;
	accessKeySecret string;
	associateTag string
	region string;
}

// Create a new AmazonRequest initialized with the given parameters
func NewRequest(accessKeyID string, accessKeySecret string, associateTag string, region string) *AmazonRequest {
	return &AmazonRequest{accessKeyID, accessKeySecret, associateTag, region}
}

// Perform an ItemLookup request.
//
// Usage:
// ids := []string{"01289328","2837423"}
// response,err := request.ItemLookup(ids, "Medium", "ASIN")
func (self AmazonRequest) ItemLookup(itemIds []string, responseGroup string, idType string) (ItemLookupResponse, error) {
	now := time.Now()
	arguments := make(map[string]string)
	arguments["AWSAccessKeyId"] = self.accessKeyID
	arguments["Version"] = "2011-08-01"
	arguments["Timestamp"] = now.Format("2006-01-02T15:04:05Z")
	arguments["Operation"] = "ItemLookup"
	arguments["Service"] = "AWSEcommerceService"
	arguments["AssociateTag"] = self.associateTag // optional
	arguments["ItemId"] = strings.Join(itemIds, ",")
	arguments["ResponseGroup"] = responseGroup
	arguments["IdType"] = idType

	// Sort the keys otherwise Amazon hash will be
	// different from mine and the request will fail
	keys := make([]string, 0, len(arguments))
	for argument := range arguments {
		keys = append(keys, argument)
	}
	sort.Strings(keys)

	// There's probably a more efficient way to concatenate strings, not a big deal though.
	var queryString string
	for _, key := range keys {
		escapedArg := url.QueryEscape(arguments[key])
		queryString += fmt.Sprintf("%s=%s", key, escapedArg)

		// Add '&' but only if it's not the the last argument
		if key != keys[len(keys)-1] {
			queryString += "&"
		}
	}

	// Hash & Sign
	var err error
	domain := service_domains[self.region]

	data := "GET\n" + domain + "\n/onca/xml\n" + queryString
	hash := hmac.New(sha256.New, []byte(self.accessKeySecret))
	hash.Write([]byte(data))
	signature := url.QueryEscape(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	queryString = fmt.Sprintf("%s&Signature=%s", queryString, signature)

	// Do request
	requestURL := fmt.Sprintf("http://%s/onca/xml?%s", domain, queryString)
	content, err := doRequest(requestURL)

	if err != nil {
		return ItemLookupResponse{}, err
	}

	return unmarshal(content)
}

// TODO add "Accept-Encoding": "gzip" and override UserAgent
// which is set to Go http package.
func doRequest(requestURL string) ([]byte, error) {
	var httpResponse *http.Response
	var err error
	var contents []byte

	httpResponse, err = http.Get(requestURL)

	if err != nil {
		return []byte(""), err
	}

	contents, err = ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()

	return contents, err
}

