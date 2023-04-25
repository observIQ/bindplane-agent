package m365receiver

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type m365Client struct {
	cfg    *Config
	client *http.Client
	root   string
	token  string
}

func newM365Client(c *http.Client, cfg *Config) *m365Client {
	return &m365Client{
		cfg:    cfg,
		root:   "https://graph.microsoft.com/v1.0/reports/",
		client: c,
	}
}

func (m *m365Client) GetCSV(endpoint string) ([]string, error) {
	req, err := http.NewRequest("GET", m.root+endpoint, nil)
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("Authorization", m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)

	//parse out 2nd line & return csv data
	_, err = csvReader.Read()
	if err != nil {
		return []string{}, err
	}
	data, err := csvReader.Read()
	if err != nil {
		return []string{}, err
	}

	return data, nil
}

// Get authorization token
type response struct {
	Token string `json:"access_token"`
}

func (m *m365Client) GetToken() error {
	auth_endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", m.cfg.Tenant_id)

	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"scope":         {"https://graph.microsoft.com/.default"},
		"client_id":     {m.cfg.Client_id},
		"client_secret": {m.cfg.Client_secret},
	}

	requestBody := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", auth_endpoint, requestBody)
	if err != nil {
		//TODO: error handling
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//TODO: error handling
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//TODO: error handling
	}

	var token response
	err = json.Unmarshal(body, &token)
	if err != nil {
		//TODO: error handling
	}

	m.token = token.Token

	return nil
}
