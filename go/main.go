package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"unicode/utf8"
)

var possibleDomains []string

type TargetEndpointRequest struct {
	Url string `json:"url"`
}

type TargetEndpointResponse struct {
	Url          string `json:"url"`
	Result       bool   `json:"result"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type TargetEndpointFormatVerification struct {
	Result       bool
	ErrorMessage string
}

func LoadPossibleDomains() []string {
	filepath := "domain_list.txt"
	dat, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer dat.Close()

	var period string = "."
	scanner := bufio.NewScanner(dat)
	for scanner.Scan() {
		line := scanner.Text()
		var appendedString = period + line
		var finalString = strings.ToLower(appendedString)
		possibleDomains = append(possibleDomains, finalString)

	}
	log.Println("Possible domains loaded in.")
	return possibleDomains
}

func VerifyEndpointUrlFormat(url string) TargetEndpointFormatVerification {

	targetEndpointFormatVerification := TargetEndpointFormatVerification{}

	if utf8.RuneCountInString(url) < 7 {
		log.Printf("Less than 7 characters in URL: %s\n", url)
		targetEndpointFormatVerification.ErrorMessage = "Target Endpoint URL does not begin with http:// or https://."
		targetEndpointFormatVerification.Result = false
		return targetEndpointFormatVerification
	}

	if utf8.RuneCountInString(url) >= 8 {
		firstEight := url[:8]
		if firstEight != "https://" {
			log.Printf("URL does not start with http:// or https://: %s\n", url)
			targetEndpointFormatVerification.ErrorMessage = "Target Endpoint URL does not begin with http:// or https://."
			targetEndpointFormatVerification.Result = false
			return targetEndpointFormatVerification
		}
	} else {
		firstSeven := url[:7]
		if firstSeven != "http://" {
			log.Printf("First 7 characters are not http:// in URL: %s\n", url)
			targetEndpointFormatVerification.ErrorMessage = "Target Endpoint URL does not begin with http:// or https://."
			targetEndpointFormatVerification.Result = false
			return targetEndpointFormatVerification
		}
	}

	i := strings.LastIndex(url, ".")

	if i < -1 {
		targetEndpointFormatVerification.ErrorMessage = "Invalid Endpoint URL domain."
		targetEndpointFormatVerification.Result = false
		return targetEndpointFormatVerification
	}

	domainEnding := url[i:]

	if !slices.Contains(possibleDomains, domainEnding) {
		targetEndpointFormatVerification.ErrorMessage = "Endpoint URL domain is not a valid domain."
		targetEndpointFormatVerification.Result = false
		return targetEndpointFormatVerification
	}

	targetEndpointFormatVerification.Result = true
	targetEndpointFormatVerification.ErrorMessage = ""

	return targetEndpointFormatVerification
}

func PingTargetEndpoint(url string) bool {
	resp, err := http.Get(url)

	if err != nil {
		return false
	}

	log.Printf("Status code for URL: %s is: %d\n", url, resp.StatusCode)

	return true
}

func CheckTargetEndpoint(w http.ResponseWriter, r *http.Request) {

	bytedata, err := io.ReadAll(r.Body)

	if err != nil {

	}

	var targetEndpointJson TargetEndpointRequest
	targetEndpointResponse := TargetEndpointResponse{}

	if err := json.Unmarshal(bytedata, &targetEndpointJson); err != nil {
		log.Println("Error occurred!")
		targetEndpointResponse.ErrorMessage = "Error with POST body."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(targetEndpointResponse)

		if err != nil {
			log.Printf("Error encoding response: %v\n", err)
		}
		return
	}
	log.Printf("Endpoint URL sent to be tested: %s\n", targetEndpointJson.Url)

	targetEndpointResponse.Url = targetEndpointJson.Url
	verifyEndpointFormatResponse := VerifyEndpointUrlFormat(targetEndpointJson.Url)

	if !verifyEndpointFormatResponse.Result {
		targetEndpointResponse.ErrorMessage = verifyEndpointFormatResponse.ErrorMessage

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(targetEndpointResponse)

		if err != nil {
			log.Printf("Error encoding response: %v\n", err)
		}
		return
	}

	targetEndpointResponse.Result = PingTargetEndpoint(targetEndpointJson.Url)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(targetEndpointResponse); err != nil {
		log.Printf("Error encoding response: %v\n", err)
	}
}

func main() {

	LoadPossibleDomains()
	http.HandleFunc("/checkTargetEndpoint", CheckTargetEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
