package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// BadNames is a name list of pods/services that shouldn't be doing any requests
// (they have invalid data).
var (
	// We discover the hostname dynamically
	Hostname, err = os.Hostname()
	BadNames      = []string{"kube-controller-manager", "kube-scheduler", "kube-scheduler-" + Hostname, "kube-controller-manager-" + Hostname}
)

// httpClient is a HTTP client with some settings for all requests.
var (
	httpClient *http.Client
)

// Start httpClient automatically.
func init() {
	httpClient = createHTTPClient()
	fmt.Println("HTTP client ready")
}

// createHTTPClient configures httpClient.
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			MaxIdleConnsPerHost: 100,
			TLSHandshakeTimeout: time.Duration(2 * time.Second),
			DisableKeepAlives:   true,
		},
		Timeout: time.Duration(3 * time.Second),
	}
	return client
}

// GetResponse returns the response from any supported request.
func GetResponse(rq *types.Request) string {
	switch rq.Action {
	case types.GET:
		fallthrough
	case types.POST:
		fallthrough
	case types.DELETE:
		return makeRequest(rq)
	default:
		err := fmt.Sprintf("Undefined request: Action %s does not exist", rq.Action)
		panic(err)
	}
}

// makeRequest processes the desired "rq" request.
func makeRequest(rq *types.Request) string {
	// Prepares the request
	req, err := http.NewRequest(rq.Action.String(), rq.URL.String(), rq.Payload)
	if err != nil {
		panic(err.Error())
	}
	// Configures the header
	req.Proto = "HTTP/1.1"
	req.Header.Set("User-Agent", "kube-nftlb-client/1.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Add(rq.Header.Key, rq.Header.Value)
	// Does the request
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	// Closes the body after reading it
	defer resp.Body.Close()
	// Reads the body from the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return string(body)
}
