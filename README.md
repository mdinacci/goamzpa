# goamzpa

A BSD licensed [Go lang](http://golang.org) library to use the _Amazon Product API_. 
Also my first `Go` project.

At the moment it supports only `ItemLookup`. Everything can change, and
probably will, use at your own peril.

## Usage
    
    package main

	import (
		"fmt"
		"amzpa"
	)

	func main() {
	    // Complete these variables with your credentials
		accessKey := "ACCESS_KEY"
		accessSecret := "ACCESS_SECRET"
		associateTag := "ASSOCIATE_TAG"
		region := "UK"
	
		request := amzpa.NewRequest(accessKey, accessSecret , associateTag, region)
		asins:= []string{"0141033576,0615314465,1470057719"}
		
		responseGroup := "Medium"
		itemsType := "ASIN"
		response,err := request.ItemLookup(asins, responseGroup, itemsType)

	    fmt.Printf("%s \n", response)
	}
 
 

## TODO
* Map the XML to a struct, so that the response is not just a big string
* ItemSearch
* Gzip compression

