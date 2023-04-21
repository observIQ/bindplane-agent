// Copyright  OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
)

const (
	typeStr = "m365"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.Stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
			CollectionInterval: 12 * time.Hour,
		},
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Timeout: 10 * time.Second,
		},
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

func createMetricsReceiver(
	_ context.Context,
	params receiver.CreateSettings,
	rConf component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	cfg := rConf.(*Config)

	//get authorization token
	token, err := getAuthorizationToken(cfg)
	if err != nil {
		return nil, err
	}

	//create receiver
	ns := newM365Scraper(params, cfg, token)
	scraper, err := scraperhelper.NewScraper(typeStr, ns.scrape, scraperhelper.WithStart(ns.start))
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(
		&cfg.ScraperControllerSettings,
		params,
		consumer,
		scraperhelper.AddScraper(scraper),
	)

}

//Getting new authentication token for metrics

type response struct {
	Token string `json:"access_token"`
}

func getAuthorizationToken(cfg *Config) (string, error) {
	auth_endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", cfg.Tenant_id)

	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"scope":         {"https://graph.microsoft.com/.default"},
		"client_id":     {cfg.Client_id},
		"client_secret": {cfg.Client_secret},
	}

	requestBody := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", auth_endpoint, requestBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var token response
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}
