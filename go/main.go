package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type TargetEndpointRequest struct {
	Url string `json:"url"`
}

type TargetEndpointResponse struct {
	Url          string `json:"url"`
	Result       bool   `json:"result"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func PingTargetEndpoint(url string) bool {
	_, err := http.Get(url)

	if err != nil {
		return false
	}
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

	targetEndpointResponse.Result = PingTargetEndpoint(targetEndpointJson.Url)
	targetEndpointResponse.Url = targetEndpointJson.Url

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(targetEndpointResponse); err != nil {
		log.Printf("Error encoding response: %v\n", err)
	}
}

func main() {
	http.HandleFunc("/checkTargetEndpoint", CheckTargetEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
