package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Telegram SDK
type Api struct {
	client  http.Client
	baseUrl url.URL
}

// Create Telegram SDK with provided token
func createApi(token string) Api {
	baseUrl, _ := url.Parse("https://api.telegram.org/bot" + token + "/")

	tr := http.Transport{}

	client := http.Client{
		Transport: &tr,
	}

	return Api{
		client:  client,
		baseUrl: *baseUrl,
	}
}

// Get information about bot itself
func (t Api) GetMe() (me Me, err error) {
	url := t.baseUrl.ResolveReference(&url.URL{Path: "getMe"})
	req, _ := http.NewRequest("GET", url.String(), nil)

	response, err := t.client.Do(req)
	if err != nil {
		return me, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&me)

	return me, err
}

// Get updates from provided offset
func (t Api) GetUpdates(offset int) (updates []Update, err error) {
	url := t.baseUrl.ResolveReference(&url.URL{Path: "getUpdates"})
	req, _ := http.NewRequest("GET", url.String()+fmt.Sprintf("?offset=%d", offset), nil)

	response, err := t.client.Do(req)
	if err != nil {
		return updates, err
	}

	defer response.Body.Close()

	var data struct {
		Ok     bool
		Result []Update
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	updates = data.Result

	return updates, err
}
