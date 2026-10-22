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

package runtime

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace/dynatrace-bindplane-otel-contrib/pkg/osinfo"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/observiq/bindplane-otel-collector/internal/extension/opampconnectionextension/internal/opamp"
	"github.com/observiq/bindplane-otel-collector/internal/extension/opampconnectionextension/internal/opamp/observiq"
	"github.com/spf13/pflag"
	"golang.org/x/net/http/httpproxy"
)

const (
	defaultConnectionTestTimeout = 30 * time.Second

	// managerPathENV is the environment variable the collector reads the manager config path from
	managerPathENV = "MANAGER_YAML_PATH"
)

var (
	// proxyEnvVars are the environment variables Go reads to select a proxy
	proxyEnvVars = []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy"}

	// noProxyEnvVars are the environment variables Go reads to bypass a proxy
	noProxyEnvVars = []string{"NO_PROXY", "no_proxy"}
)

var (
	// connectionTestMaxRetries is the number of times a handshake rejected with 429 Too Many Requests is retried
	connectionTestMaxRetries = 3

	// connectionTestRetryInterval is the initial delay between retries when the server does not send Retry-After
	connectionTestRetryInterval = time.Second

	// connectionTestSleep waits between retries, and is replaced in tests
	connectionTestSleep = time.Sleep
)

// loadServiceProxyEnv is overridden in tests so they do not depend on the host's service configuration
var loadServiceProxyEnv = connectionTest.loadServiceProxyEnv

// connectionTest reports the progress of a connection test
type connectionTest struct {
	out io.Writer
}

func (c connectionTest) printf(format string, a ...any) {
	_, _ = fmt.Fprintf(c.out, format+"\n", a...)
}

// RunConnectionTest performs a websocket handshake with the OpAMP endpoint in the manager config,
// using the same TLS, proxy, and headers as the agent. No OpAMP messages are sent.
// Progress is written to out.
func RunConnectionTest(args []string, out io.Writer) error {
	c := connectionTest{out: out}

	defaultPaths := defaultManagerPaths()
	flags := pflag.NewFlagSet("test-connection", pflag.ContinueOnError)
	flags.SetOutput(out)
	managerPath := flags.String("manager", "", "Path to the manager config. Defaults to the first that exists of: "+strings.Join(defaultPaths, ", "))
	timeout := flags.Duration("timeout", defaultConnectionTestTimeout, "Maximum time to wait for the handshake")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return nil
		}
		return err
	}
	if *timeout <= 0 {
		return errors.New("--timeout must be greater than zero")
	}

	// Use the proxy environment of the collector service when it can be determined
	if err := loadServiceProxyEnv(c); err != nil {
		return err
	}
	setProxyVars := c.reportProxyEnv()

	if *managerPath == "" {
		path, err := findManagerConfig(defaultPaths)
		if err != nil {
			return err
		}
		*managerPath = path
	}

	cfg, err := loadManagerConfig(*managerPath)
	if err != nil {
		return err
	}
	c.printf("Loaded manager config from %s", *managerPath)

	endpoint, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint %q: %w", cfg.Endpoint, err)
	}
	if endpoint.Scheme != "ws" && endpoint.Scheme != "wss" {
		return fmt.Errorf("unsupported endpoint scheme %q: must be ws or wss", endpoint.Scheme)
	}

	agentID := cfg.AgentID
	if agentID == opamp.EmptyAgentID {
		u, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("new uuidv7: %w", err)
		}
		agentID = opamp.AgentIDFromUUID(u)
		c.printf("Manager config has no agent_id, using temporary agent ID %s (not saved)", agentID)
	} else {
		c.printf("Agent ID: %s", agentID)
	}
	if cfg.SecretKey == nil {
		c.printf("Secret key: not set")
	} else {
		c.printf("Secret key: set")
	}

	tlsCfg, err := cfg.ToTLS(nil)
	if err != nil {
		return fmt.Errorf("failed creating TLS config: %w", err)
	}
	c.reportTLS(cfg, tlsCfg, endpoint)

	// Mirror the dialer the opamp-go websocket client builds
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = tlsCfg

	proxyURL, err := c.envProxy(endpoint, setProxyVars)
	if err != nil {
		return err
	}
	dialer.Proxy = func(*http.Request) (*url.URL, error) {
		return proxyURL, nil
	}

	hostname, err := osinfo.Hostname()
	if err != nil {
		c.printf("Warning: failed to retrieve hostname: %s", err)
	}
	header := observiq.ConnectionHeaders(http.Header{}, cfg.GetSecretKey(), agentID, hostname)

	var (
		conn      *websocket.Conn
		resp      *http.Response
		serverTLS *tls.ConnectionState
	)
	retryInterval := connectionTestRetryInterval
	for retry := 1; ; retry++ {
		c.printf("Connecting to %s", endpoint)
		conn, resp, serverTLS, err = handshake(&dialer, endpoint, header, *timeout)
		if err == nil || resp == nil || resp.StatusCode != http.StatusTooManyRequests {
			break
		}
		if retry > connectionTestMaxRetries {
			c.printf("Server responded with %s%s, giving up after %d retries", resp.Status, bindplaneReason(resp), connectionTestMaxRetries)
			break
		}

		// Bindplane rejects connections with 429 before the websocket upgrade, such as when too many
		// agents are connecting or it detects a duplicate agent. Retry-After is honored when set,
		// otherwise the delay doubles with each retry.
		wait, source := retryInterval, "backoff"
		if retryAfter, ok := parseRetryAfter(resp); ok {
			wait, source = retryAfter, "Retry-After"
		}
		c.printf("Server responded with %s%s, waiting %s (%s) before retry %d of %d",
			resp.Status, bindplaneReason(resp), wait, source, retry, connectionTestMaxRetries)
		connectionTestSleep(wait)
		retryInterval *= 2
	}
	if err != nil {
		c.reportServerTLS(serverTLS, err)
		return handshakeError(err, resp, proxyURL)
	}
	defer conn.Close()

	c.printf("Handshake succeeded: %s", resp.Status)
	c.reportServerTLS(serverTLS, nil)

	_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
	return nil
}

// handshake performs a single websocket handshake. The server's TLS state is returned so it can be
// reported even if the handshake fails.
func handshake(dialer *websocket.Dialer, endpoint *url.URL, header http.Header, timeout time.Duration) (*websocket.Conn, *http.Response, *tls.ConnectionState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// With an https proxy, the proxy's TLS handshake is reported first, so the last one is the server's
	var serverTLS *tls.ConnectionState
	if endpoint.Scheme == "wss" {
		ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
			TLSHandshakeDone: func(state tls.ConnectionState, _ error) {
				serverTLS = &state
			},
		})
	}

	conn, resp, err := dialer.DialContext(ctx, endpoint.String(), header)
	return conn, resp, serverTLS, err
}

// parseRetryAfter returns the delay requested by the Retry-After header, given in seconds or as an HTTP date
func parseRetryAfter(resp *http.Response) (time.Duration, bool) {
	value := resp.Header.Get("Retry-After")
	if value == "" {
		return 0, false
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil && seconds >= 0 {
		// Clamp before converting, since multiplying a large value by time.Second overflows to a negative duration
		return time.Duration(min(seconds, int64(math.MaxInt64/time.Second))) * time.Second, true
	}
	if t, err := http.ParseTime(value); err == nil {
		return max(time.Until(t), 0), true
	}
	return 0, false
}

// bindplaneReason returns the rejection reason Bindplane sends in the X-Bindplane-Reason header, formatted for output
func bindplaneReason(resp *http.Response) string {
	if reason := resp.Header.Get("X-Bindplane-Reason"); reason != "" {
		return " (" + reason + ")"
	}
	return ""
}

// defaultManagerPaths returns the manager config paths to try, in order, when --manager is not set
func defaultManagerPaths() []string {
	paths := []string{"./manager.yaml"}
	switch runtime.GOOS {
	case "linux":
		paths = append(paths, "/opt/observiq-otel-collector/manager.yaml")
	case "windows":
		if programFiles := os.Getenv("ProgramFiles"); programFiles != "" {
			paths = append(paths, filepath.Join(programFiles, "observIQ OpenTelemetry Collector", "manager.yaml"))
		}
	}
	if path := os.Getenv(managerPathENV); path != "" {
		paths = append(paths, path)
	}
	return paths
}

// findManagerConfig returns the first of paths that exists
func findManagerConfig(paths []string) (string, error) {
	for _, path := range paths {
		_, err := os.Stat(path)
		switch {
		case err == nil:
			return path, nil
		case !errors.Is(err, os.ErrNotExist):
			return "", fmt.Errorf("failed to read manager config: %w", err)
		}
	}
	return "", fmt.Errorf("manager config not found at %s, use --manager to set the path", strings.Join(paths, ", "))
}

// loadManagerConfig reads the manager config without bootstrapping or modifying it
func loadManagerConfig(path string) (*opamp.Config, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("manager config not found at %s", path)
		}
		return nil, fmt.Errorf("failed to read manager config: %w", err)
	}

	cfg, err := opamp.ParseConfig(path)
	if err != nil {
		return nil, err
	}

	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("manager config %s has no endpoint", path)
	}

	return cfg, nil
}

// reportProxyEnv prints the proxy environment variables and returns the names of the proxy variables that are set
func (c connectionTest) reportProxyEnv() []string {
	var setProxyVars []string
	for _, kv := range os.Environ() {
		name, value, _ := strings.Cut(kv, "=")
		switch {
		case value == "":
		case slices.Contains(proxyEnvVars, name):
			setProxyVars = append(setProxyVars, name)
			c.printf("Proxy environment variable %s is set: %s", name, redactProxy(value))
		case slices.Contains(noProxyEnvVars, name):
			c.printf("Proxy environment variable %s is set: %s", name, value)
		}
	}
	return setProxyVars
}

// reportTLS prints the TLS settings and switches the endpoint to wss when TLS is configured, as the agent does
func (c connectionTest) reportTLS(cfg *opamp.Config, tlsCfg *tls.Config, endpoint *url.URL) {
	if tlsCfg == nil {
		if endpoint.Scheme == "wss" {
			c.printf("TLS: using the system certificate store")
		}
		return
	}

	if endpoint.Scheme == "ws" {
		endpoint.Scheme = "wss"
		c.printf("tls_config is set, using wss instead of ws (the agent does the same)")
	}
	if cfg.TLS.CAFile != nil {
		c.printf("TLS: using CA file %s (replaces the system certificate store)", *cfg.TLS.CAFile)
	} else {
		c.printf("TLS: using the system certificate store")
	}
	if cfg.TLS.CertFile != nil {
		c.printf("TLS: using client certificate %s", *cfg.TLS.CertFile)
	}
	if tlsCfg.InsecureSkipVerify {
		c.printf("TLS: insecure_skip_verify is enabled, the server certificate will not be verified")
	}
}

// reportServerTLS prints the negotiated TLS version and the server certificate. A TLS-inspecting
// proxy presents a certificate issued by its own CA instead of the server's.
func (c connectionTest) reportServerTLS(state *tls.ConnectionState, err error) {
	var certs []*x509.Certificate
	if state != nil {
		if state.HandshakeComplete {
			c.printf("TLS version: %s", tls.VersionName(state.Version))
		}
		certs = state.PeerCertificates
	}

	// Certificates that fail verification are only available from the error
	var certErr *tls.CertificateVerificationError
	if len(certs) == 0 && errors.As(err, &certErr) {
		certs = certErr.UnverifiedCertificates
	}

	if len(certs) > 0 {
		c.printf("Server certificate subject: %s", certs[0].Subject)
		c.printf("Server certificate issuer: %s", certs[0].Issuer)
	}
}

// envProxy returns the proxy selected from the environment for the endpoint, or nil if no proxy applies
func (c connectionTest) envProxy(endpoint *url.URL, setProxyVars []string) (*url.URL, error) {
	// gorilla/websocket selects the proxy using http and https in place of ws and wss
	target := &url.URL{Scheme: "http", Host: endpoint.Host}
	if endpoint.Scheme == "wss" {
		target.Scheme = "https"
	}

	proxyURL, err := httpproxy.FromEnvironment().ProxyFunc()(target)
	if err != nil {
		// The error contains the raw proxy value, which may include credentials
		return nil, errors.New("invalid proxy environment variable")
	}

	switch {
	case proxyURL != nil:
		var source string
		for _, name := range proxyEnvVarsFor(endpoint) {
			if os.Getenv(name) != "" {
				source = name
				break
			}
		}
		c.printf("Using proxy %s from %s", proxyURL.Redacted(), source)
	case len(setProxyVars) > 0:
		c.printf("No proxy applies to %s (no %s set, or the host matches NO_PROXY or is a loopback address)", endpoint.Host, proxyEnvVarsFor(endpoint)[0])
	default:
		c.printf("No proxy configured")
	}
	return proxyURL, nil
}

// proxyEnvVarsFor returns the proxy environment variables that apply to the endpoint's scheme
func proxyEnvVarsFor(endpoint *url.URL) []string {
	if endpoint.Scheme == "wss" {
		return []string{"HTTPS_PROXY", "https_proxy"}
	}
	return []string{"HTTP_PROXY", "http_proxy"}
}

// redactProxy returns the proxy value with any password removed
func redactProxy(proxy string) string {
	u, err := url.Parse(proxy)
	if err != nil || u.Scheme == "" || u.Host == "" {
		u, err = url.Parse("http://" + proxy)
		if err != nil {
			return "<invalid proxy URL>"
		}
	}
	return u.Redacted()
}

// handshakeError adds a hint for the likely cause of a failed handshake
func handshakeError(err error, resp *http.Response, proxyURL *url.URL) error {
	if resp != nil {
		switch {
		case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
			return fmt.Errorf("server responded with %s, check secret_key in the manager config: %w", resp.Status, err)
		case resp.StatusCode == http.StatusNotFound:
			return fmt.Errorf("server responded with %s, check the endpoint path in the manager config: %w", resp.Status, err)
		case resp.StatusCode >= 300 && resp.StatusCode < 400:
			return fmt.Errorf("server responded with %s redirecting to %q, consider updating the endpoint in the manager config: %w", resp.Status, resp.Header.Get("Location"), err)
		case resp.StatusCode == http.StatusTooManyRequests:
			msg := fmt.Sprintf("server responded with %s%s", resp.Status, bindplaneReason(resp))
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				msg += ", Retry-After: " + retryAfter
			}
			return fmt.Errorf("%s. Bindplane rejects connections when too many agents are connecting or when it detects "+
				"another connection with the same agent ID: %w", msg, err)
		case proxyURL != nil:
			return fmt.Errorf("server responded with %s through proxy %s, the proxy may not support websockets. "+
				"TLS-inspecting proxies often remove the Upgrade header; if the server certificate issuer is not the server's CA, "+
				"ask for the OpAMP endpoint to be excluded from inspection: %w", resp.Status, proxyURL.Redacted(), err)
		default:
			return fmt.Errorf("server responded with %s: %w", resp.Status, err)
		}
	}

	var dnsErr *net.DNSError
	var certErr *tls.CertificateVerificationError
	var authorityErr x509.UnknownAuthorityError
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()):
		return fmt.Errorf("timed out connecting to the server: %w", err)
	case errors.As(err, &dnsErr):
		return fmt.Errorf("failed to resolve %s: %w", dnsErr.Name, err)
	case errors.As(err, &certErr) || errors.As(err, &authorityErr):
		return fmt.Errorf("server certificate verification failed, set tls_config.ca_file to the CA that signed the server certificate (or the CA of a TLS-inspecting proxy): %w", err)
	case proxyURL != nil:
		return fmt.Errorf("connection through proxy %s failed: %w", proxyURL.Redacted(), err)
	default:
		return err
	}
}
