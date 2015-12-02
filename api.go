package goconsul

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

var (
	endpointRegister string
	endpointQuery    string
	client           *http.Client
)

// NewClient returns a client struct with
// a preset timeout.
func NewClient() *http.Client {
	return &http.Client{}
}

// SetData adds a value to the consul key-value store.
func SetData(key, value string) (err error) {
	endpointStore := fmt.Sprintf("http://localhost:%s/v1/kv/%s", port, key)
	r, err := http.NewRequest("PUT", endpointStore, bytes.NewBufferString(string(value)))
	if err != nil {
		return err
	}
	log.Printf("adding data for key %s for %s\n", key, consul.Name)
	client = NewClient()
	_, err = client.Do(r)
	if err != nil {
		return err
	}
	return nil
}

// GetData retrieves a value from the consul key-value store.
func GetData(key string) (value string, err error) {
	endpointStore := fmt.Sprintf("http://localhost:%s/v1/kv/%s", port, key)
	r, err := http.NewRequest("GET", endpointStore, nil)
	if err != nil {
		return "", err
	}
	log.Printf("retrieving data for key %s for %s\n", key, consul.Name)
	client = NewClient()
	_, err = client.Do(r)
	if err != nil {
		return "", err
	}
	return value, nil
}

// doRegistration registers a service
// with the local consul agent.
func doRegistration(data []byte) (err error) {
	// Register service via PUT request.
	endpointRegister = fmt.Sprintf("http://localhost:%s/v1/agent/service/register", port)
	r, err := http.NewRequest("PUT", endpointRegister, bytes.NewBufferString(string(data)))

	if err != nil {
		return err
	}

	log.Printf("sending service registration request for %s\n", consul.Name)

	client = NewClient()
	_, err = client.Do(r)

	if err != nil {
		return err
	}
	return nil
}

// queryService queries the local consul agent for a service.
func queryService() (resp *http.Response, err error) {
	// Query service via HTTP API for confirmation in console.
	endpointQuery = fmt.Sprintf("http://localhost:%s/v1/catalog/service/%s", port, consul.Name)
	req, err := http.NewRequest("GET", endpointQuery, nil)
	if err != nil {
		return nil, err
	}
	client = NewClient()
	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
