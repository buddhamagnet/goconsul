package goconsul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	endpointRegister string
	endpointQuery    string
	client           *http.Client
)

type KV struct {
	Value []byte `json:"value"`
}

// NewClient returns a client struct with
// a preset timeout.
func NewClient() *http.Client {
	return &http.Client{}
}

// SetData adds a value to the consul key-value store.
func SetData(key string, value []byte) (err error) {
	endpointStore := fmt.Sprintf("http://localhost:%s/v1/kv/%s", port, key)
	r, err := http.NewRequest("PUT", endpointStore, bytes.NewBuffer(value))
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
func GetData(key string) (string, error) {
	var vals []KV
	endpointStore := fmt.Sprintf("http://localhost:%s/v1/kv/%s", port, key)
	r, err := http.NewRequest("GET", endpointStore, nil)
	if err != nil {
		return "", err
	}
	log.Printf("retrieving data for key %s for %s\n", key, consul.Name)
	client = NewClient()
<<<<<<< HEAD
	_, err = client.Do(r)
	if err != nil {
		return "", err
	}
	return value, nil
=======
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(data, &vals)
	return string(vals[0].Value), nil
>>>>>>> ece8ee509c8469b7c57ff7ed3cbd9158150b179a
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
