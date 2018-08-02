package outlyer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIResponse represents an Outlyer API response when HTTP response
type APIResponse struct {
	Code        int
	Body        []byte
	ErrorDetail error
}

// Get will set the API token and default headers before issuing a GET request to Outlyer API
func Get(endpoint string) ([]byte, error) {
	baseURL := UserConfig.GetString("api-url")
	completeURL := baseURL + endpoint

	client := http.Client{}
	req, err := http.NewRequest("GET", completeURL, nil)
	if err != nil {
		return nil, err
	}

	// Add request headers
	token := UserConfig.GetString("api-token")
	req.Header.Add(http.CanonicalHeaderKey("Authorization"), fmt.Sprintf("Bearer %s", token))
	req.Header.Add(http.CanonicalHeaderKey("Content-Type"), UserConfig.GetString("headers.post.content-type"))
	commonHeaders := UserConfig.GetStringMapString("headers.common")
	for k, v := range commonHeaders {
		req.Header.Add(http.CanonicalHeaderKey(k), v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s", content)
	}

	return content, nil
}

// Patch will set the API token and default headers before issuing a PATCH request to Outlyer API
func Patch(endpoint string, payload []byte) (*APIResponse, error) {
	return send(endpoint, "PATCH", payload)
}

// send wil issue an HTTP request for the given Outlyer API endpoint with the method and payload provided
func send(endpoint, method string, payload []byte) (*APIResponse, error) {
	baseURL := UserConfig.GetString("api-url")
	completeURL := baseURL + endpoint

	client := http.Client{}
	req, err := http.NewRequest(method, completeURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// Add request headers
	token := UserConfig.GetString("api-token")
	req.Header.Add(http.CanonicalHeaderKey("Authorization"), fmt.Sprintf("Bearer %s", token))
	req.Header.Add(http.CanonicalHeaderKey("Content-Type"), UserConfig.GetString("headers.post.content-type"))
	commonHeaders := UserConfig.GetStringMapString("headers.common")
	for k, v := range commonHeaders {
		req.Header.Add(http.CanonicalHeaderKey(k), v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &APIResponse{Code: resp.StatusCode, Body: content, ErrorDetail: getHTTPErrorBy(resp.StatusCode)}, nil
}

func getHTTPErrorBy(responseCode int) error {
	var err error
	if responseCode >= 400 && responseCode < 500 {
		err = fmt.Errorf("client error")
	}
	if responseCode == 400 {
		err = fmt.Errorf("incorrect resource definition, ensure all fields are correct and try again")
	}
	if responseCode == 401 {
		err = fmt.Errorf("you don't have permissions to perform this operation")
	}
	if responseCode == 404 {
		err = fmt.Errorf("resource not found")
	}
	if responseCode >= 500 && responseCode < 600 {
		err = fmt.Errorf("Outlyer API is unavailable, try again later")
	}
	return err
}
