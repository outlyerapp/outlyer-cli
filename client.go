package outlyer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// Get will perform a GET request to the Outlyer API
func Get(endpoint string) ([]byte, error) {
	baseURL := UserConfig.GetString("api-url")
	url := baseURL + endpoint

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add request headers
	token := UserConfig.GetString("api-token")
	req.Header.Add(http.CanonicalHeaderKey("Authorization"), fmt.Sprintf("Bearer %s", token))
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
		viper.SetConfigType("json")
		viper.ReadConfig(bytes.NewReader(content))
		error := fmt.Errorf("Error Code: %s\nError Message: %s", viper.GetString("status"), viper.GetString("detail"))
		return nil, error
	}

	return content, nil
}
