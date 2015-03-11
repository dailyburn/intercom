package intercom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const cBaseUrl = "https://api.intercom.io/"

const DEFAULT_MAX_RESULTS = 200

const MPOST = "POST"
const MGET = "GET"

// Intercom is a client object with functions to make reuqests to the Intercom api
type Intercom struct {
	client     *http.Client
	baseurl    string
	auth       Auth
	maxResults int
}

// Auth contains username and password attributes used for api request authentication
type Auth struct {
	AppId  string
	ApiKey string
}

// NewIntercomClient returns an instance of the Intercom api client
func NewIntercomClient(appId, apiKey string, maxResults int) *Intercom {
	if maxResults == -1 {
		maxResults = DEFAULT_MAX_RESULTS
	}
	d := &Intercom{client: &http.Client{}, auth: Auth{appId, apiKey}, maxResults: maxResults}

	return d
}

// UpdateUser posts updated user data and returns an error if the POST fails
func (i *Intercom) UpdateUser(params map[string]interface{}) error {
	urlStr := i.buildUrl("users", nil)
	_, err := i.execRequest(MPOST, urlStr, params)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

// CreateEvent posts event data and returns an error if the POST fails
func (i *Intercom) CreateEvent(params map[string]interface{}) error {
	urlStr := i.buildUrl("events", nil)
	_, err := i.execRequest(MPOST, urlStr, params)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

// buildUrl creates a url for the given path and url parameters
func (i *Intercom) buildUrl(path string, params map[string]string) string {
	var aUrl *url.URL
	aUrl, err := url.Parse(cBaseUrl)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	aUrl.Path += path
	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	aUrl.RawQuery = parameters.Encode()
	return aUrl.String()
}

// execRequest executes an arbitrary request for the given method and url returning the contents of the response in a map, or an error
func (i *Intercom) execRequest(method, aUrl string, params map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("EXEC REQUEST: ", aUrl, " method:", method, " params:", params)

	// json string encode the params for the POST body if there are any
	var body io.Reader
	if params != nil && method == MPOST {
		b, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Json error: ", err)
		}
		body = bytes.NewBuffer(b)
		fmt.Println("BODY: ", string(b))
	}

	req, err := http.NewRequest(method, aUrl, body)
	if err != nil {
		fmt.Println("execRequest error: ", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(i.auth.AppId, i.auth.ApiKey)

	fmt.Println("URL: ", req.URL)

	resp, rerr := i.client.Do(req)
	if rerr != nil {
		fmt.Println("req error: ", rerr)
		return nil, rerr
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		data, derr := ioutil.ReadAll(resp.Body)
		if derr != nil {
			fmt.Println("Error reading response: ", derr)
			return nil, derr
		}
		fmt.Println("DATA: ", string(data))
		var parsed map[string]interface{}
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			return nil, err
		}
		return parsed, nil
	case 202:
		return nil, nil
	case 404:
		return nil, errors.New("not-found")
	case 429:
		resets_at, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)
		return nil, RateLimitError(resets_at)
	case 500, 502, 503, 504:
		return nil, errors.New("server")
	default:
		return nil, errors.New("unknown, error code: " + fmt.Sprintf("%d", resp.StatusCode))
	}

}
