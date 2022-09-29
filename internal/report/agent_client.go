// Copyright  observIQ, Inc.
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

package report

import (
	"crypto/tls"
	"net/http"
)

var _ Client = (*AgentClient)(nil)

// AgentClient is a basic client that injects agent specific information in request headers
type AgentClient struct {
	agentID   string
	secretKey *string

	client *http.Client
}

// NewAgentClient creates a new AgentClient
func NewAgentClient(agentID string, secretKey *string, tlsConfig *tls.Config) *AgentClient {
	return &AgentClient{
		agentID:   agentID,
		secretKey: secretKey,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}
}

// Do injects agent specific information into headers then sends the request
func (a *AgentClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Agent-ID", a.agentID)

	if a.secretKey != nil {
		req.Header.Add("X-BindPlane-Secret-Key", *a.secretKey)
	}

	return a.client.Do(req)
}
