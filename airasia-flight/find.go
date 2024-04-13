package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func main() {
	start := time.Now()
	end := start.AddDate(1, 0, 0)
	wg := sync.WaitGroup{}

	fmt.Println("search until: ", end)

	// SG
	fmt.Println("SG")
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		// fmt.Println(d.Format("02-01-2006"))

		date := d.Format("02-01-2006")
		wg.Add(1)
		go func() {
			doSG(date)
			wg.Done()
			time.Sleep(1 * time.Second)
		}()
	}

	wg.Wait()

	// JKT
	fmt.Println("JKT")
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		// fmt.Println(d.Format("02-01-2006"))

		date := d.Format("02-01-2006")
		wg.Add(1)
		go func() {
			doJKT(date)
			wg.Done()
			time.Sleep(1 * time.Second)
		}()
	}

	wg.Wait()
}

type AirAsiaResponse struct {
	Notifications []struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Threshold string `json:"threshold"`
	} `json:"notifications"`
}

func doSG(date string) {
	client := &http.Client{}
	var data = strings.NewReader(`{"consumerId":"ListingAPP","flightJourney":{"journeyType":"O","journeyDetails":[{"departDate":"` + date + `","origin":"SIN","destination":"JKT"}],"passengers":{"adult":1,"child":0,"infant":0}},"searchContext":{"promocode":"DLB-aEbZo8-8b586","sort":"cheapest","filters":{"cabin":{"applyMixedClasses":false,"cabinClass":"ECONOMY"},"carriers":{"allowAllCarriers":false,"excludedCarriers":[],"onlyAllowedCarriers":["AK","Z2","D7","FD","QZ","XJ"]},"duration":{"maxStopoverTimeInHrs":0,"maxTravelTimeInHrs":0,"minStopoverTimeInHrs":0},"stops":{"allowOvernight":false,"stopType":"ANY"}}},"userContext":{"currency":"IDR","geoId":"SG","locale":"en-gb","platform":"web","experimentVariants":[]}}`)
	req, err := http.NewRequest("POST", "https://flights.airasia.com/web/fp/search/flights/v5/aggregated-results?type=paired&isPromoMessagesByCode=true&include_list=searchResults,content,currency,vouchers,upsellSnap,upsellPremiumFlatBed&airlineProfile=all&isOriginCity=true&isDestinationCity=true&uce=false&page=1", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "Bearer ") // add your token here
	req.Header.Set("channel_hash", "c5e9028b4295dcf4d7c239af8231823b520c3cc15b99ab04cde71d0ab18d65bc")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("ga-id", "GA1.2.1463702013.1709708241")
	req.Header.Set("origin", "https://www.airasia.com")
	req.Header.Set("referer", "https://www.airasia.com/")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-type", "MEMBER")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := AirAsiaResponse{}
	err = jsoniter.Unmarshal(bodyText, &response)
	if err != nil {
		log.Fatal(err)
	}

	if len(response.Notifications) > 0 {
		// fmt.Println(date+" response.Notifications[0].Threshold: ", response.Notifications[0].Threshold)
		if response.Notifications[0].Threshold == "SUCCESS" {
			fmt.Println("found date: ", date)
		}
	}
}

func doJKT(date string) {
	client := &http.Client{}
	var data = strings.NewReader(`{"consumerId":"ListingAPP","flightJourney":{"journeyType":"O","journeyDetails":[{"departDate":"` + date + `","origin":"JKT","destination":"SIN"}],"passengers":{"adult":1,"child":0,"infant":0}},"searchContext":{"promocode":"DLB-aEbZo8-8b586","sort":"cheapest","filters":{"cabin":{"applyMixedClasses":false,"cabinClass":"ECONOMY"},"carriers":{"allowAllCarriers":false,"excludedCarriers":[],"onlyAllowedCarriers":["AK","Z2","D7","FD","QZ","XJ"]},"duration":{"maxStopoverTimeInHrs":0,"maxTravelTimeInHrs":0,"minStopoverTimeInHrs":0},"stops":{"allowOvernight":false,"stopType":"ANY"}}},"userContext":{"currency":"IDR","geoId":"SG","locale":"en-gb","platform":"web","experimentVariants":[]}}`)
	req, err := http.NewRequest("POST", "https://flights.airasia.com/web/fp/search/flights/v5/aggregated-results?type=paired&isPromoMessagesByCode=true&include_list=searchResults,content,currency,vouchers,upsellSnap,upsellPremiumFlatBed&airlineProfile=all&isOriginCity=true&isDestinationCity=true&uce=false&page=1", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "Bearer ") // add your token here
	req.Header.Set("channel_hash", "c5e9028b4295dcf4d7c239af8231823b520c3cc15b99ab04cde71d0ab18d65bc")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("ga-id", "GA1.2.1463702013.1709708241")
	req.Header.Set("origin", "https://www.airasia.com")
	req.Header.Set("referer", "https://www.airasia.com/")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-type", "MEMBER")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := AirAsiaResponse{}
	err = jsoniter.Unmarshal(bodyText, &response)
	if err != nil {
		log.Fatal(err)
	}

	if len(response.Notifications) > 0 {
		// fmt.Println(date+" response.Notifications[0].Threshold: ", response.Notifications[0].Threshold)
		if response.Notifications[0].Threshold == "SUCCESS" {
			fmt.Println("found date: ", date)
		}
	}
}
