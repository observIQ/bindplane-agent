package report

import (
	"net/http"
)

var _ Client = (*AgentClient)(nil)

// AgentClient is a basic client that inject agent specific information in request headers
type AgentClient struct {
	agentID   string
	secretKey *string

	client *http.Client
}

// NewAgentClient creates a new AgentClient
func NewAgentClient(agentID string, secretKey *string) *AgentClient {
	return &AgentClient{
		agentID:   agentID,
		secretKey: secretKey,
		client:    http.DefaultClient,
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
