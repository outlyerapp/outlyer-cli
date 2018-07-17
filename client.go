package outlyer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// Get will set up the API token and default headers before issuing a GET request to the Outlyer API
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
	commonHeaders := UserConfig.GetStringMapString("headers.common")
	for k, v := range commonHeaders {
		req.Header.Add(http.CanonicalHeaderKey(k), v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error communicating with Outlyer API", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		viper.SetConfigType("yaml")
		viper.ReadConfig(bytes.NewReader(content))
		error := fmt.Errorf("\nError Code: %s\nError Message: %s\n", viper.GetString("status"), viper.GetString("detail"))
		return nil, error
	}

	return content, nil
}
