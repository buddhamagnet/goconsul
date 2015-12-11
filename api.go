package goconsul

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
)

var (
	endpointRegister string
	endpointQuery    string
	consulClient     *api.Client
	webClient        *http.Client
)

type KV struct {
	Value []byte `json:"value"`
}

func init() {
	consulClient, _ = api.NewClient(api.DefaultConfig())
	webClient = new(http.Client)
}

// SetData adds a value to the consul key-value store.
func SetData(key string, value []byte) (err error) {
	kv := consulClient.KV()
	_, err = kv.Put(&api.KVPair{Key: key, Value: value}, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetData retrieves a value from the consul key-value store.
func GetData(key string) ([]byte, error) {
	var data []byte
	kv := consulClient.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return data, err
	}
	return pair.Value, nil
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

	_, err = webClient.Do(r)

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
	resp, err = webClient.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
