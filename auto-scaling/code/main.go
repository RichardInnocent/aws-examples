package main

import (
	"io"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Get)
	http.ListenAndServe(":8080", nil)
}

func Get(responseWriter http.ResponseWriter, _ *http.Request) {
	instanceID, err := GetInstanceID()
	if err != nil {
		responseWriter.WriteHeader(500)
		responseWriter.Write([]byte(err.Error()))
		return
	}
	responseWriter.Write([]byte(fmt.Sprintf("Hello from %s\n", instanceID))
}

func GetInstanceID() (string, error) {
	response, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}
	instanceIDBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}
	return string(instanceIDBytes), nil
}
