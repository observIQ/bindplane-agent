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
	"bytes"
	"encoding/pem"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	goruntime "runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

const testAgentIDString = "01HX2DWEQZ045KQR3VG0EYEZ94"

// clearProxyEnv ensures proxy variables and services from the test environment do not affect the test
func clearProxyEnv(t *testing.T) {
	for _, name := range append(proxyEnvVars, noProxyEnvVars...) {
		t.Setenv(name, "")
	}

	original := loadServiceProxyEnv
	loadServiceProxyEnv = func(connectionTest) error { return nil }
	t.Cleanup(func() { loadServiceProxyEnv = original })
}

func writeManagerConfig(t *testing.T, contents string) string {
	path := filepath.Join(t.TempDir(), "manager.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0600))
	return path
}

// newOpAMPServer starts a websocket server that records the handshake headers
func newOpAMPServer(t *testing.T, tlsServer bool) (*httptest.Server, func() http.Header) {
	var mu sync.Mutex
	var gotHeader http.Header
	upgrader := websocket.Upgrader{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotHeader = r.Header.Clone()
		mu.Unlock()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		_, _, _ = conn.ReadMessage()
	})

	var srv *httptest.Server
	if tlsServer {
		srv = httptest.NewTLSServer(handler)
	} else {
		srv = httptest.NewServer(handler)
	}
	t.Cleanup(srv.Close)

	return srv, func() http.Header {
		mu.Lock()
		defer mu.Unlock()
		return gotHeader
	}
}

type connectProxy struct {
	mu       sync.Mutex
	connects int
}

func (p *connectProxy) connectCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connects
}

// newConnectProxy starts an HTTP CONNECT proxy that tunnels every request to backend
func newConnectProxy(t *testing.T, backend string) (*httptest.Server, *connectProxy) {
	p := &connectProxy{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		p.mu.Lock()
		p.connects++
		p.mu.Unlock()

		backendConn, err := net.Dial("tcp", backend)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		clientConn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			backendConn.Close()
			return
		}
		_, _ = clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

		go func() {
			defer backendConn.Close()
			_, _ = io.Copy(backendConn, clientConn)
		}()
		go func() {
			defer clientConn.Close()
			_, _ = io.Copy(clientConn, backendConn)
		}()
	}))
	t.Cleanup(srv.Close)
	return srv, p
}

func wsEndpoint(srv *httptest.Server) string {
	return "ws://" + srv.Listener.Addr().String() + "/v1/opamp"
}

func wssEndpoint(srv *httptest.Server) string {
	return "wss://" + srv.Listener.Addr().String() + "/v1/opamp"
}

func writeCAFile(t *testing.T, srv *httptest.Server) string {
	path := filepath.Join(t.TempDir(), "ca.crt")
	data := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})
	require.NoError(t, os.WriteFile(path, data, 0600))
	return path
}

func TestRunConnectionTest_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		env         map[string]string
		expectedErr string
	}{
		{
			name:        "invalid timeout",
			args:        []string{"--timeout", "0s"},
			expectedErr: "--timeout must be greater than zero",
		},
		{
			name:        "unknown flag",
			args:        []string{"--bogus"},
			expectedErr: "unknown flag: --bogus",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearProxyEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			out := &bytes.Buffer{}
			err := RunConnectionTest(tc.args, out)
			require.ErrorContains(t, err, tc.expectedErr)
			require.NotContains(t, out.String(), "secret")
		})
	}
}

func TestRunConnectionTest_Help(t *testing.T) {
	clearProxyEnv(t)
	out := &bytes.Buffer{}
	require.NoError(t, RunConnectionTest([]string{"--help"}, out))
	require.Contains(t, out.String(), "--timeout")
}

func TestRunConnectionTest_MissingManagerConfig(t *testing.T) {
	clearProxyEnv(t)
	path := filepath.Join(t.TempDir(), "manager.yaml")

	err := RunConnectionTest([]string{"--manager", path}, &bytes.Buffer{})
	require.ErrorContains(t, err, "manager config not found at "+path)

	_, err = os.Stat(path)
	require.ErrorIs(t, err, os.ErrNotExist, "manager config must not be created")
}

func TestDefaultManagerPaths(t *testing.T) {
	t.Setenv(managerPathENV, "/env/manager.yaml")

	paths := defaultManagerPaths()
	require.Equal(t, "./manager.yaml", paths[0])
	require.Equal(t, "/env/manager.yaml", paths[len(paths)-1])
	switch goruntime.GOOS {
	case "linux":
		require.Equal(t, []string{"./manager.yaml", "/opt/observiq-otel-collector/manager.yaml", "/env/manager.yaml"}, paths)
	case "windows":
		require.Contains(t, paths, filepath.Join(os.Getenv("ProgramFiles"), "observIQ OpenTelemetry Collector", "manager.yaml"))
	}
}

func TestFindManagerConfig(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.yaml")
	second := filepath.Join(dir, "second.yaml")
	missing := filepath.Join(dir, "missing.yaml")
	require.NoError(t, os.WriteFile(first, []byte("endpoint: ws://localhost\n"), 0600))
	require.NoError(t, os.WriteFile(second, []byte("endpoint: ws://localhost\n"), 0600))

	path, err := findManagerConfig([]string{missing, first, second})
	require.NoError(t, err)
	require.Equal(t, first, path)

	_, err = findManagerConfig([]string{missing})
	require.ErrorContains(t, err, "manager config not found at "+missing+", use --manager to set the path")
}

func TestRunConnectionTest_DefaultManagerConfig(t *testing.T) {
	clearProxyEnv(t)
	srv, _ := newOpAMPServer(t, false)
	path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))
	t.Chdir(t.TempDir())

	t.Run("environment variable", func(t *testing.T) {
		if _, err := os.Stat("/opt/observiq-otel-collector/manager.yaml"); goruntime.GOOS == "linux" && err == nil {
			t.Skip("the default Linux install path takes precedence over the environment variable")
		}
		t.Setenv(managerPathENV, path)
		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest(nil, out))
		require.Contains(t, out.String(), "Loaded manager config from "+path)
	})

	t.Run("current directory before environment variable", func(t *testing.T) {
		t.Setenv(managerPathENV, path)
		contents, err := os.ReadFile(path)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile("manager.yaml", contents, 0600))
		t.Cleanup(func() { _ = os.Remove("manager.yaml") })

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest(nil, out))
		require.Contains(t, out.String(), "Loaded manager config from ./manager.yaml")
	})
}

func TestRunConnectionTest_InvalidManagerConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		expectedErr string
	}{
		{
			name:        "missing endpoint",
			config:      "agent_id: " + testAgentIDString + "\n",
			expectedErr: "has no endpoint",
		},
		{
			name:        "unsupported scheme",
			config:      "endpoint: http://localhost:3001/v1/opamp\n",
			expectedErr: `unsupported endpoint scheme "http"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearProxyEnv(t)
			path := writeManagerConfig(t, tc.config)
			err := RunConnectionTest([]string{"--manager", path}, &bytes.Buffer{})
			require.ErrorContains(t, err, tc.expectedErr)
		})
	}
}

func TestRunConnectionTest_Success(t *testing.T) {
	clearProxyEnv(t)
	srv, gotHeader := newOpAMPServer(t, false)
	path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nsecret_key: my-secret\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

	out := &bytes.Buffer{}
	require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))

	header := gotHeader()
	require.Equal(t, "Secret-Key my-secret", header.Get("Authorization"))
	require.Equal(t, testAgentIDString, header.Get("Agent-ID"))
	require.True(t, strings.HasPrefix(header.Get("User-Agent"), "observiq-otel-collector/"))

	output := out.String()
	require.Contains(t, output, "Agent ID: "+testAgentIDString)
	require.Contains(t, output, "Secret key: set")
	require.Contains(t, output, "No proxy configured")
	require.Contains(t, output, "Handshake succeeded: 101 Switching Protocols")
	require.NotContains(t, output, "my-secret")
}

func TestRunConnectionTest_TemporaryAgentID(t *testing.T) {
	clearProxyEnv(t)
	srv, gotHeader := newOpAMPServer(t, false)
	path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\n", wsEndpoint(srv)))

	out := &bytes.Buffer{}
	require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
	require.Contains(t, out.String(), "using temporary agent ID")
	require.Contains(t, out.String(), "Secret key: not set")
	require.NotEmpty(t, gotHeader().Get("Agent-ID"))

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("endpoint: %s\n", wsEndpoint(srv)), string(contents), "manager config must not be modified")
}

func TestRunConnectionTest_ServerRejects(t *testing.T) {
	testCases := []struct {
		name        string
		status      int
		expectedErr string
	}{
		{
			name:        "unauthorized",
			status:      http.StatusUnauthorized,
			expectedErr: "server responded with 401 Unauthorized, check secret_key",
		},
		{
			name:        "not found",
			status:      http.StatusNotFound,
			expectedErr: "server responded with 404 Not Found, check the endpoint path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearProxyEnv(t)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			}))
			t.Cleanup(srv.Close)
			path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

			err := RunConnectionTest([]string{"--manager", path}, &bytes.Buffer{})
			require.ErrorContains(t, err, tc.expectedErr)
		})
	}
}

func TestRunConnectionTest_TLS(t *testing.T) {
	t.Run("ca file", func(t *testing.T) {
		clearProxyEnv(t)
		srv, _ := newOpAMPServer(t, true)
		caFile := writeCAFile(t, srv)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\ntls_config:\n  ca_file: %s\n", wssEndpoint(srv), testAgentIDString, caFile))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Contains(t, out.String(), "TLS: using CA file "+caFile)
		require.Contains(t, out.String(), "TLS version: TLS 1.3")
		require.Contains(t, out.String(), "Server certificate issuer:")
	})

	t.Run("tls config switches ws to wss", func(t *testing.T) {
		clearProxyEnv(t)
		srv, _ := newOpAMPServer(t, true)
		caFile := writeCAFile(t, srv)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\ntls_config:\n  ca_file: %s\n", wsEndpoint(srv), testAgentIDString, caFile))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Contains(t, out.String(), "tls_config is set, using wss instead of ws")
		require.Contains(t, out.String(), "Connecting to wss://")
	})

	t.Run("untrusted certificate", func(t *testing.T) {
		clearProxyEnv(t)
		srv, _ := newOpAMPServer(t, true)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wssEndpoint(srv), testAgentIDString))

		out := &bytes.Buffer{}
		err := RunConnectionTest([]string{"--manager", path}, out)
		require.ErrorContains(t, err, "server certificate verification failed, set tls_config.ca_file")
		require.Contains(t, out.String(), "TLS: using the system certificate store")
		require.Contains(t, out.String(), "Server certificate issuer: O=Acme Co")
	})

	t.Run("server certificate reported when handshake is rejected", func(t *testing.T) {
		clearProxyEnv(t)
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		t.Cleanup(srv.Close)
		caFile := writeCAFile(t, srv)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\ntls_config:\n  ca_file: %s\n", wssEndpoint(srv), testAgentIDString, caFile))

		out := &bytes.Buffer{}
		err := RunConnectionTest([]string{"--manager", path}, out)
		require.ErrorContains(t, err, "server responded with 400 Bad Request: websocket: bad handshake")
		require.NotContains(t, err.Error(), "proxy")
		require.Contains(t, out.String(), "TLS version: TLS 1.3")
		require.Contains(t, out.String(), "Server certificate issuer: O=Acme Co")
	})
}

func TestRunConnectionTest_ProxyEnv(t *testing.T) {
	t.Run("proxy applies", func(t *testing.T) {
		clearProxyEnv(t)
		srv, _ := newOpAMPServer(t, false)
		proxySrv, proxy := newConnectProxy(t, srv.Listener.Addr().String())
		proxyEnv := "http://user:secret@" + proxySrv.Listener.Addr().String()
		t.Setenv("HTTP_PROXY", proxyEnv)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: ws://bindplane.example.com/v1/opamp\nagent_id: %s\n", testAgentIDString))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Equal(t, 1, proxy.connectCount())

		output := out.String()
		redacted := "http://user:xxxxx@" + proxySrv.Listener.Addr().String()
		require.Contains(t, output, "Proxy environment variable HTTP_PROXY is set: "+redacted)
		require.Contains(t, output, "Using proxy "+redacted+" from HTTP_PROXY")
		require.NotContains(t, output, "secret")
	})

	t.Run("no proxy match", func(t *testing.T) {
		clearProxyEnv(t)
		srv, _ := newOpAMPServer(t, false)
		proxySrv, proxy := newConnectProxy(t, srv.Listener.Addr().String())
		t.Setenv("HTTP_PROXY", proxySrv.URL)
		t.Setenv("NO_PROXY", "bindplane.example.com")
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Equal(t, 0, proxy.connectCount())
		require.Contains(t, out.String(), "Proxy environment variable NO_PROXY is set: bindplane.example.com")
		require.Contains(t, out.String(), "No proxy applies to "+srv.Listener.Addr().String())
	})

	t.Run("proxy does not support websockets", func(t *testing.T) {
		clearProxyEnv(t)
		// The server sees a request without the Upgrade header, as when a TLS-inspecting proxy removes it
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		t.Cleanup(srv.Close)
		proxySrv, proxy := newConnectProxy(t, srv.Listener.Addr().String())
		t.Setenv("HTTPS_PROXY", proxySrv.URL)
		caFile := writeCAFile(t, srv)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: wss://example.com/v1/opamp\nagent_id: %s\ntls_config:\n  ca_file: %s\n", testAgentIDString, caFile))

		out := &bytes.Buffer{}
		err := RunConnectionTest([]string{"--manager", path}, out)
		require.Equal(t, 1, proxy.connectCount())
		require.ErrorContains(t, err, "server responded with 400 Bad Request through proxy "+proxySrv.URL+", the proxy may not support websockets")
		require.Contains(t, out.String(), "Server certificate issuer: O=Acme Co")
	})

	t.Run("scheme mismatch", func(t *testing.T) {
		clearProxyEnv(t)
		srv, _ := newOpAMPServer(t, false)
		proxySrv, proxy := newConnectProxy(t, srv.Listener.Addr().String())
		t.Setenv("HTTPS_PROXY", proxySrv.URL)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Equal(t, 0, proxy.connectCount())
		require.Contains(t, out.String(), "no HTTP_PROXY set")
	})
}

// recordRetryWaits replaces the retry sleep so tests do not wait, and returns the waits that were requested
func recordRetryWaits(t *testing.T) func() []time.Duration {
	var mu sync.Mutex
	var waits []time.Duration
	original := connectionTestSleep
	connectionTestSleep = func(d time.Duration) {
		mu.Lock()
		defer mu.Unlock()
		waits = append(waits, d)
	}
	t.Cleanup(func() { connectionTestSleep = original })

	return func() []time.Duration {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(waits)
	}
}

// newRejectingServer returns 429 for the first rejections requests, then upgrades to a websocket.
// A negative rejections value rejects every request.
func newRejectingServer(t *testing.T, rejections int, header map[string]string) (*httptest.Server, func() int) {
	var mu sync.Mutex
	attempts := 0
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attempts++
		reject := rejections < 0 || attempts <= rejections
		mu.Unlock()

		if reject {
			for k, v := range header {
				w.Header().Set(k, v)
			}
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		_, _, _ = conn.ReadMessage()
	}))
	t.Cleanup(srv.Close)

	return srv, func() int {
		mu.Lock()
		defer mu.Unlock()
		return attempts
	}
}

func TestRunConnectionTest_TooManyRequests(t *testing.T) {
	t.Run("backoff without Retry-After", func(t *testing.T) {
		clearProxyEnv(t)
		waits := recordRetryWaits(t)
		srv, attempts := newRejectingServer(t, 2, nil)
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Equal(t, 3, attempts())
		require.Equal(t, []time.Duration{time.Second, 2 * time.Second}, waits())
		require.Contains(t, out.String(), "Server responded with 429 Too Many Requests, waiting 1s (backoff) before retry 1 of 3")
		require.Contains(t, out.String(), "Server responded with 429 Too Many Requests, waiting 2s (backoff) before retry 2 of 3")
		require.Contains(t, out.String(), "Handshake succeeded: 101 Switching Protocols")
	})

	t.Run("honors Retry-After", func(t *testing.T) {
		clearProxyEnv(t)
		waits := recordRetryWaits(t)
		srv, attempts := newRejectingServer(t, 1, map[string]string{"Retry-After": "30", "X-Bindplane-Reason": "duplicate agent detected"})
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

		out := &bytes.Buffer{}
		require.NoError(t, RunConnectionTest([]string{"--manager", path}, out))
		require.Equal(t, 2, attempts())
		require.Equal(t, []time.Duration{30 * time.Second}, waits())
		require.Contains(t, out.String(), "Server responded with 429 Too Many Requests (duplicate agent detected), waiting 30s (Retry-After) before retry 1 of 3")
	})

	t.Run("gives up after three retries", func(t *testing.T) {
		clearProxyEnv(t)
		waits := recordRetryWaits(t)
		srv, attempts := newRejectingServer(t, -1, map[string]string{"Retry-After": "5", "X-Bindplane-Reason": "duplicate agent detected"})
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

		out := &bytes.Buffer{}
		err := RunConnectionTest([]string{"--manager", path}, out)
		require.Equal(t, 4, attempts())
		require.Equal(t, []time.Duration{5 * time.Second, 5 * time.Second, 5 * time.Second}, waits())
		require.Contains(t, out.String(), "before retry 3 of 3")
		require.NotContains(t, out.String(), "retry 4")
		require.Contains(t, out.String(), "Server responded with 429 Too Many Requests (duplicate agent detected), giving up after 3 retries")
		require.ErrorContains(t, err, "server responded with 429 Too Many Requests (duplicate agent detected), Retry-After: 5. Bindplane rejects connections")
	})

	t.Run("Retry-After that overflows a duration", func(t *testing.T) {
		clearProxyEnv(t)
		waits := recordRetryWaits(t)
		srv, attempts := newRejectingServer(t, -1, map[string]string{"Retry-After": "99999999999999999"})
		path := writeManagerConfig(t, fmt.Sprintf("endpoint: %s\nagent_id: %s\n", wsEndpoint(srv), testAgentIDString))

		err := RunConnectionTest([]string{"--manager", path}, &bytes.Buffer{})
		require.ErrorContains(t, err, "Retry-After: 99999999999999999")
		require.Equal(t, 4, attempts())
		for _, wait := range waits() {
			require.Greater(t, wait, 100*365*24*time.Hour, "a large Retry-After must not overflow to a short or negative wait")
		}
	})
}

func TestParseRetryAfter(t *testing.T) {
	maxSeconds := time.Duration(math.MaxInt64/time.Second) * time.Second
	testCases := []struct {
		name       string
		value      string
		expectedOK bool
		check      func(t *testing.T, d time.Duration)
	}{
		{name: "missing", value: ""},
		{name: "invalid", value: "soon"},
		{name: "negative", value: "-1"},
		{
			name:       "seconds",
			value:      "30",
			expectedOK: true,
			check:      func(t *testing.T, d time.Duration) { require.Equal(t, 30*time.Second, d) },
		},
		{
			name:       "seconds that overflow a duration",
			value:      "99999999999999999",
			expectedOK: true,
			check:      func(t *testing.T, d time.Duration) { require.Equal(t, maxSeconds, d) },
		},
		{
			name:       "http date",
			value:      time.Now().Add(time.Minute).UTC().Format(http.TimeFormat),
			expectedOK: true,
			check: func(t *testing.T, d time.Duration) {
				require.Greater(t, d, 50*time.Second)
				require.LessOrEqual(t, d, time.Minute)
			},
		},
		{
			name:       "http date far in the future",
			value:      "Fri, 31 Dec 9999 23:59:59 GMT",
			expectedOK: true,
			check:      func(t *testing.T, d time.Duration) { require.Greater(t, d, 100*365*24*time.Hour) },
		},
		{
			name:       "http date in the past",
			value:      time.Now().Add(-time.Minute).UTC().Format(http.TimeFormat),
			expectedOK: true,
			check:      func(t *testing.T, d time.Duration) { require.Zero(t, d) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tc.value != "" {
				resp.Header.Set("Retry-After", tc.value)
			}

			d, ok := parseRetryAfter(resp)
			require.Equal(t, tc.expectedOK, ok)
			if tc.check != nil {
				tc.check(t, d)
			}
		})
	}
}
