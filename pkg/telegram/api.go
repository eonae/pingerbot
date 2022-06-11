package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"pingerbot/pkg/helpers"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// Telegram SDK
type Api struct {
	client  http.Client
	baseUrl url.URL
}

// Create Telegram SDK with provided token
func NewApi(token string) Api {
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

type tgRequest struct {
	Method  string
	Path    string
	Query   map[string]string
	Body    interface{}
	Headers map[string]string
}

type tgResponse struct {
	Ok          bool        `json:"ok"`
	ErrorCode   int         `json:"error_code,omitempty"`
	Description string      `json:"description"`
	Result      interface{} `json:"result"`
}

const ResponseTypeMismatch = helpers.Error("Response type mismatch")

// Get information about bot itself
func (t Api) GetMe() (me Me, err error) {
	err = t.sendRequest(tgRequest{Method: "GET", Path: "getMe"}, &me)
	if err != nil {
		return
	}

	return
}

// Get updates from provided offset
func (t Api) GetUpdates(offset int64, timeout time.Duration) (updates []Update, err error) {
	err = t.sendRequest(tgRequest{
		Method: "GET",
		Path:   "getUpdates",
		Query: map[string]string{
			"offset":  strconv.Itoa(int(offset)),
			"timeout": strconv.Itoa(int(timeout.Milliseconds())),
		},
	}, &updates)

	return
}

func (t Api) SendMessage(command OutgoingMessage) (message IncomingMessage, err error) {
	err = t.sendRequest(tgRequest{
		Method: "POST",
		Path:   "sendMessage",
		Body:   command,
	}, &message)

	return
}

func (t Api) createRequest(opts tgRequest) (req *http.Request, err error) {
	url := t.baseUrl.ResolveReference(&url.URL{Path: opts.Path}).String()

	if opts.Query != nil {
		i := 0
		for k, v := range opts.Query {
			c := "&"
			if i == 0 {
				c = "?"
			}

			url += fmt.Sprintf("%s%s=%s", c, k, v)
			i++
		}
	}

	var body io.Reader

	if opts.Body != nil {
		buf, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, err
		}

		logrus.Debug("Sending body (raw):", string(buf))

		body = bytes.NewReader(buf)
	}

	req, err = http.NewRequest(opts.Method, url, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

func (t Api) sendRequest(opts tgRequest, result interface{}) (err error) {
	req, err := t.createRequest(opts)
	if err != nil {
		return
	}

	res, err := t.client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	logrus.Debug("Raw data", string(buf))

	var data tgResponse

	// Crutial thing. If not doing this, data.Result
	// will be unmarshalled to map[string]interface{}
	data.Result = result

	err = json.Unmarshal(buf, &data)
	if err != nil {
		return
	}

	if data.Ok {
		return nil
	}

	return fmt.Errorf("%d %s", data.ErrorCode, data.Description)
}
