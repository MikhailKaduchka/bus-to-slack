package main

// Getting departures from api.9292.nl and send them to slack by calling incoming webhook
// Uses LOCATION_URL and SLACK_WEBHOOK variables.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type LocationDepartures struct {
	Location struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		StopType string `json:"stopType"`
		Name     string `json:"name"`
		Place    struct {
			Name        string `json:"name"`
			RegionCode  string `json:"regionCode"`
			RegionName  string `json:"regionName"`
			ShowRegion  bool   `json:"showRegion"`
			CountryCode string `json:"countryCode"`
			CountryName string `json:"countryName"`
			ShowCountry bool   `json:"showCountry"`
		} `json:"place"`
		LatLong struct {
			Lat  float64 `json:"lat"`
			Long float64 `json:"long"`
		} `json:"latLong"`
		Urls struct {
			NlNL string `json:"nl-NL"`
			EnGB string `json:"en-GB"`
		} `json:"urls"`
	} `json:"location"`
	Tabs []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Locations []struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			StopType string `json:"stopType"`
			Name     string `json:"name"`
			Place    struct {
				Name        string `json:"name"`
				RegionCode  string `json:"regionCode"`
				RegionName  string `json:"regionName"`
				ShowRegion  bool   `json:"showRegion"`
				CountryCode string `json:"countryCode"`
				CountryName string `json:"countryName"`
				ShowCountry bool   `json:"showCountry"`
			} `json:"place"`
			LatLong struct {
				Lat  float64 `json:"lat"`
				Long float64 `json:"long"`
			} `json:"latLong"`
			Urls struct {
				NlNL string `json:"nl-NL"`
				EnGB string `json:"en-GB"`
			} `json:"urls"`
		} `json:"locations"`
		Departures []struct {
			Time            string      `json:"time"`
			DestinationName string      `json:"destinationName"`
			ViaNames        interface{} `json:"viaNames"`
			Mode            struct {
				Type string `json:"type"`
				Name string `json:"name"`
			} `json:"mode"`
			OperatorName    string      `json:"operatorName"`
			Service         string      `json:"service"`
			Platform        interface{} `json:"platform"`
			PlatformChanged bool        `json:"platformChanged"`
			Remark          interface{} `json:"remark"`
			RealtimeState   string      `json:"realtimeState"`
			RealtimeText    string      `json:"realtimeText"`
		} `json:"departures"`
	} `json:"tabs"`
}

type SlackMessage struct {
	Text string `json:"text"`
}

func getDataByURL(url string) []byte {
	httpClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func convertJsonToLocation(jsonBody []byte) LocationDepartures {
	location := LocationDepartures{}
	jsonErr := json.Unmarshal(jsonBody, &location)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return location
}

func getMessageTextFromLocation(location LocationDepartures) string {
	message := ""

	for i := 0; i < len(location.Tabs[0].Departures); i++ {
		message += "At: " + location.Tabs[0].Departures[i].Time + " To: " + location.Tabs[0].Departures[i].DestinationName + "\n"
	}

	// log.Println("Message to send: " + message)
	return message
}

func sendToSlack(url string, message SlackMessage) {
	httpClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(message)

	req, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.StatusCode >= 400 {
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		log.Fatal("Send to Slack failed with status: " + res.Status + " Body: " + string(body))
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Handler(request Request) (Response, error) {
	url := getEnv("LOCATION_URL", "https://api.9292.nl/0.1/locations/STOP_NAME/departure-times?lang=en-GB")
	webhook := getEnv("SLACK_WEBHOOK", "https://hooks.slack.com/services/TOKEN-HERE")

	log.Println("Make a request to API...")
	response := getDataByURL(url)

	log.Println("Converting response to location object...")
	location := convertJsonToLocation(response)

	log.Println("Making a message for sending to Slack")
	message := SlackMessage{getMessageTextFromLocation(location)}

	log.Println("Sending a message to Slack...")
	sendToSlack(webhook, message)

	return Response{
		Message: fmt.Sprintf("Processed request ID %f", request.ID),
		Ok:      true,
	}, nil
}

func main() {
	log.Println("Started")
	lambda.Start(Handler)
	log.Println("Finished")
}
