package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

const (
	eveDefaultEndpoint   = "https://localhost"
	eveEndpointEnvVar    = "EVECTL_ENDPOINT"
	eveUserNameEnvVar    = "EVECTL_USERNAME"
	evePasswordEnvVar    = "EVECTL_PASSWORD"
	eveCAFileEnvVar      = "EVECTL_CA_FILE"
	eveTLSNoVerifyEnvVar = "EVECTL_TLS_NOVERIFY"
)

// ErrorAuth is the error returned if a 401 is returned by an API request.
var ErrorAuth = fmt.Errorf("authentication failed")

// ErrorNotFound is the error returned if a 404 is returned by an API request.
var ErrorNotFound = fmt.Errorf("resource not found")

// Client represents a single connection to Eve
type Client struct {
	Url           *url.URL
	Username      string
	Password      string
	HttpClient    *http.Client
	DefaultHeader http.Header
}

// NewDefaultClient returns a client with the default behavior
func NewDefaultClient() *Client {
	eveEndpoint := os.Getenv(eveEndpointEnvVar)
	if eveEndpoint == "" {
		eveEndpoint = eveDefaultEndpoint
	}

	client, err := NewClient(eveEndpoint)
	if err != nil {
		panic(err)
	}

	return client
}

// NewClient returns a client to the specified Eve endpoint
func NewClient(urlString string) (*Client, error) {
	if len(urlString) == 0 {
		return nil, fmt.Errorf("client: missing url")
	}

	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	username := os.Getenv(eveUserNameEnvVar)
	if username != "" {
		log.Debugf("Using EVE_USERNAME (%s)", username)
	}
	password := os.Getenv(evePasswordEnvVar)

	client := &Client{
		Url:           parsedUrl,
		Username:      username,
		Password:      password,
		DefaultHeader: make(http.Header),
	}

	userAgent := fmt.Sprintf("EveGo/1.0 (%s)", runtime.Version())
	client.DefaultHeader.Set("User-Agent", userAgent)
	client.DefaultHeader.Set("Content-Type", "application/json")
	client.DefaultHeader.Set("Accept", "application/json")

	if err := client.init(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) init() error {
	tlsConfig := &tls.Config{}

	if os.Getenv(eveTLSNoVerifyEnvVar) != "" {
		tlsConfig.InsecureSkipVerify = true
	}

	caFile := os.Getenv(eveCAFileEnvVar)
	if caFile != "" {
		log.Debugf("Loading CA File: %s", caFile)
		pool, err := c.loadCAFile(caFile)
		if err != nil {
			log.Error(err)
		}

		tlsConfig.RootCAs = pool
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	c.HttpClient = &http.Client{Transport: transport}

	return nil
}

func (c *Client) loadCAFile(caFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading CA File: %s", err)
	}

	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, fmt.Errorf("Error loading CA File: Couldn't parse PEM in: %s", caFile)
	}

	return pool, nil
}

// RequestInput represents data to build a request object
type RequestInput struct {
	Params     map[string]string
	Headers    map[string]string
	Body       io.Reader
	BodyLength int64
}

// Request creates a new HTTP request using the supplied parameters
func (c *Client) Request(verb, resourcePath string, input *RequestInput) (*http.Request, error) {
	u := *c.Url
	u.Path = path.Join(c.Url.Path, resourcePath)
	return c.rawRequest(verb, &u, input)
}

func (c *Client) rawRequest(verb string, url *url.URL, input *RequestInput) (*http.Request, error) {
	if verb == "" {
		return nil, fmt.Errorf("Client: missing verb")
	}

	if url == nil {
		return nil, fmt.Errorf("Client: missing url.URL")
	}

	if input == nil {
		return nil, fmt.Errorf("Client: missing RequestInput")
	}

	// Create the request object
	request, err := http.NewRequest(verb, url.String(), input.Body)
	if err != nil {
		return nil, err
	}

	// Add basic auth
	request.SetBasicAuth(c.Username, c.Password)

	// Set our default headers first
	for k, v := range c.DefaultHeader {
		request.Header[k] = v
	}

	// Add any request headers
	for k, v := range input.Headers {
		request.Header.Add(k, v)
	}

	// Add content-length if we have it
	if input.BodyLength > 0 {
		request.ContentLength = input.BodyLength
	}

	log.Debugf("Raw request: %#v", request)

	return request, nil
}

// checkResponse is a response wrapper that can handle the different status
// codes that can returned by eve
func checkResponse(resp *http.Response, err error) (*http.Response, error) {
	// If we already have an error, just return right away
	if err != nil {
		return resp, err
	}

	log.Debugf("Response: %d (%s)", resp.StatusCode, resp.Status)
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		log.Error("Response: error copying response body")
	} else {
		log.Debugf("Response: %s", buf.String())
		resp.Body.Close()
		resp.Body = &bytesReadCloser{&buf}
	}

	switch resp.StatusCode {
	case 200:
		return resp, nil
	case 401:
		return nil, ErrorAuth
	case 404:
		return nil, ErrorNotFound
	default:
		return nil, fmt.Errorf("Client: %s", resp.Status)
	}
}

// decodeJson is used to decode a JSON body into an interface.
func decodeJson(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	return d.Decode(out)
}

// bytesReadCloser is a simple wrapper around a bytes buffer that implements
// Close as a noop.
type bytesReadCloser struct {
	*bytes.Buffer
}

func (nrc *bytesReadCloser) Close() error {
	// we don't actually have to do anything here, since the buffer is just some
	// data in memory and the error is initialized to no-error
	return nil
}
