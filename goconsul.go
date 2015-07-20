package goconsul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	endpointRegister string
	endpointQuery    string
	client           *http.Client
	port             string
)

type Consul struct {
	Name string
}

func init() {
	port = os.Getenv("CONSUL_PORT")
	if port == "" {
		port = "8500"
	}
	registerService()
}

// NewClient returns a client struct with
// a preset timeout.
func NewClient() *http.Client {
	return &http.Client{}
}

// registerService registers the service outlined in goconsul.json
// with the local consul agent.
func registerService() {
	var consul Consul
	config, err := os.Open("goconsul.json")

	if err != nil {
		log.Fatal(err)
	}
	defer config.Close()

	data, err := ioutil.ReadAll(config)

	if err != nil {
		log.Fatal(err)
	}

	_ = json.Unmarshal(data, &consul)

	endpointRegister = fmt.Sprintf("http://localhost:%s/v1/agent/service/register", port)
	r, err := http.NewRequest("PUT", endpointRegister, bytes.NewBufferString(string(data)))

	log.Printf("sending service registration request for %s\n", consul.Name)

	client := NewClient()
	_, err = client.Do(r)

	if err != nil {
		log.Fatal(err)
	}

	endpointQuery = fmt.Sprintf("http://localhost:%s/v1/catalog/service/%s", port, consul.Name)
	r, err = http.NewRequest("GET", endpointQuery, nil)
	resp, err := client.Do(r)
	if resp.StatusCode == 200 {
		response, _ := ioutil.ReadAll(resp.Body)
		log.Println("service registration complete, please check response:")
		log.Println(string(response))
	} else {
		log.Fatalf("service registration failed: %d\n", resp.StatusCode)
	}
}
