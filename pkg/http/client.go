package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zevenet/kube-nftlb/pkg/env"
	"github.com/zevenet/kube-nftlb/pkg/types"
)

var (
	// HTTP client settings
	httpClient = &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			MaxIdleConnsPerHost: 100,
			TLSHandshakeTimeout: time.Duration(2 * time.Second),
			DisableKeepAlives:   true,
		},
		Timeout: time.Duration(3 * time.Second),
	}

	// nftlb settings
	key      = env.GetString("NFTLB_KEY")
	protocol = env.GetString("NFTLB_PROTOCOL")
	host     = env.GetString("NFTLB_HOST")
	port     = env.GetInt("NFTLB_PORT")
)

// Send returns the response from a request.
func Send(requestData *types.RequestData) ([]byte, error) {
	// Prepare the request
	request, err := http.NewRequest(requestData.Method, types.URL(protocol, host, port, requestData.Path), requestData.Body)
	if err != nil {
		return nil, err
	}

	// Set key:value pairs in header
	request.Proto = "HTTP/1.1"
	request.Header.Set("User-Agent", "kube-nftlb-client/1.0")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Key", key)

	// Do the request and get response
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	// Close body after reading it
	defer response.Body.Close()

	// Read body from the response
	return ioutil.ReadAll(response.Body)
}
