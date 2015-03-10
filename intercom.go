package intercom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const cBaseUrl = "https://api.intercom.io/"

const DEFAULT_MAX_RESULTS = 200

const MPOST = "POST"
const MGET = "GET"

// Jira is a client object with functions to make reuqests to the jira api
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

// NewJiraClient returns an instance of the Jira api client
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
	data, err := i.execRequest(MPOST, urlStr, params)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Println("DATA: ", string(data))
	return nil
}

// CreateEvent posts event data and returns an error if the POST fails
func (i *Intercom) CreateEvent(params map[string]interface{}) error {
	urlStr := i.buildUrl("events", nil)
	data, err := i.execRequest(MPOST, urlStr, params)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Println("DATA: ", string(data))
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

// execRequest executes an arbitrary request for the given method and url returning the contents of the response in []byte or an error
func (i *Intercom) execRequest(method, aUrl string, params map[string]interface{}) ([]byte, error) {
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

	data, derr := ioutil.ReadAll(resp.Body)
	if derr != nil {
		fmt.Println("Error reading response: ", derr)
		return nil, derr
	}

	return data, nil
}
