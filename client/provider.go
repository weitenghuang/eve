package client

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve"
)

// GetProvider retrieves a provider with the given name
func (c *Client) GetProvider(name string) *eve.Provider {
	endpoint := fmt.Sprintf("/provider/%s", name)
	input := &RequestInput{
		Params:     make(map[string]string),
		Headers:    make(map[string]string),
		Body:       nil,
		BodyLength: 0,
	}
	req, err := c.Request("GET", endpoint, input)
	if err != nil {
		log.Fatalf("GetProvider: %s", err)
	}

	resp, err := checkResponse(c.HttpClient.Do(req))
	if err != nil {
		log.Fatalf("GetProvider: %s", err)
	}

	var provider *eve.Provider
	if err := decodeJson(resp, &provider); err != nil {
		log.Fatalf("GetProvider: %s", err)
	}

	return provider
}
