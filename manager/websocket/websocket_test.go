package websocket

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/observiq/observiq-collector/manager/message"
	"github.com/stretchr/testify/require"
)

func TestOpenSuccess(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	conn.Close()
}

func TestOpenFailure(t *testing.T) {
	server := NewServer(t)
	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.Error(t, err)
	require.Nil(t, conn)
}

func TestHandleReceive(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	errChan := make(chan error, 1)
	in := make(chan *message.Message, 5)
	go func() { errChan <- HandleReceive(ctx, conn, in) }()

	sentMessage, err := message.New("test", &map[string]interface{}{})
	require.NoError(t, err)

	server.out <- sentMessage
	receivedMessage := <-in
	require.Equal(t, sentMessage, receivedMessage)
	require.Equal(t, 0, len(errChan))
}

func TestHandleReceiveCtx(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := make(chan *message.Message, 5)
	err = HandleReceive(ctx, conn, in)
	require.Error(t, err)
	require.Contains(t, err.Error(), context.Canceled.Error())
}

func TestHandleReceivedClosedConn(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	in := make(chan *message.Message, 5)
	go func() { errChan <- HandleReceive(ctx, conn, in) }()
	Close(conn)

	err = <-errChan
	require.Error(t, err)
	require.Contains(t, err.Error(), ErrConnectionClosed.Error())
}

func TestHandleSend(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	out := make(chan *message.Message, 5)
	go func() { errChan <- HandleSend(ctx, conn, out) }()

	sentMessage, err := message.New("test", &map[string]interface{}{})
	require.NoError(t, err)
	out <- sentMessage

	receivedMessage := <-server.in
	require.Equal(t, sentMessage, receivedMessage)
	require.Equal(t, 0, len(errChan))
}

func TestHandleSendCtx(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	out := make(chan *message.Message, 5)
	err = HandleSend(ctx, conn, out)
	require.Error(t, err)
	require.Contains(t, err.Error(), context.Canceled.Error())
}

func TestHandleSendClosedConn(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	out := make(chan *message.Message, 5)
	go func() { errChan <- HandleSend(ctx, conn, out) }()
	Close(conn)

	sentMessage, err := message.New("test", nil)
	require.NoError(t, err)

	out <- sentMessage
	err = <-errChan
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown write error")
}

func TestHandleSendClosedChan(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan *message.Message, 5)
	close(out)

	err = HandleSend(ctx, conn, out)
	require.NoError(t, err)
}

func TestHandleTraffic(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	in := make(chan *message.Message, 5)
	out := make(chan *message.Message, 5)
	go func() { errChan <- HandleTraffic(ctx, conn, in, out) }()

	testMessage, err := message.New("test", &map[string]interface{}{})
	require.NoError(t, err)

	out <- testMessage
	receivedMessage := <-server.in
	require.Equal(t, testMessage, receivedMessage)

	server.out <- testMessage
	receivedMessage = <-in
	require.Equal(t, testMessage, receivedMessage)

	cancel()
	err = <-errChan
	require.Error(t, err)
	require.Contains(t, err.Error(), context.Canceled.Error())
}

func TestWriteBadMessage(t *testing.T) {
	server := NewServer(t)
	server.Start()
	defer server.Stop()

	conn, err := Open(context.Background(), server.WebsocketAddress(), nil)
	require.NoError(t, err)
	defer conn.Close()

	invalidContent := &map[string]interface{}{
		"chan": make(chan int),
	}
	testMessage, err := message.New("test", invalidContent)
	require.NoError(t, err)

	err = write(testMessage, conn)
	require.Error(t, err)
	require.Contains(t, err.Error(), "json: unsupported type")
}

// Server is a server used for testing.
type Server struct {
	Upgrader   websocket.Upgrader
	address    string
	port       int
	httpServer *http.Server
	in         chan *message.Message
	out        chan *message.Message
}

// HandleRequest converts an http request to a websocket connection and handles incoming traffic.
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_ = HandleTraffic(context.Background(), conn, s.in, s.out)
}

// TCPAddress returns the TCP address of the server.
func (s *Server) TCPAddress() string {
	return fmt.Sprintf("%s:%d", s.address, s.port)
}

// WebsocketAddress returns the websocket address of the server.
func (s *Server) WebsocketAddress() string {
	return fmt.Sprintf("ws://%s:%d", s.address, s.port)
}

// Start will initiate the server and begin listening for requests.
func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/", s.HandleRequest)
	s.httpServer = &http.Server{
		Handler:      router,
		Addr:         s.TCPAddress(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	ready := make(chan struct{})
	go func() {
		addr := s.httpServer.Addr
		if addr == "" {
			addr = ":http"
		}

		ln, _ := net.Listen("tcp", addr)

		ready <- struct{}{}

		_ = s.httpServer.Serve(ln)
	}()
	<-ready
}

// Stop will close the server and stop listening for requests.
func (s *Server) Stop() {
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}

// NewServer returns a new mock server.
func NewServer(t *testing.T) *Server {
	port, err := findOpenPort()
	require.NoError(t, err, "Could not find open port for test")
	fmt.Printf("Websocket Server using port %d\n", port)

	return &Server{
		Upgrader: websocket.Upgrader{},
		address:  "127.0.0.1",
		port:     port,
		in:       make(chan *message.Message, 10),
		out:      make(chan *message.Message, 10),
	}
}

// findOpenPort attempts to find an open port on the localhost.
func findOpenPort() (int, error) {
	for i := 0; i < 10; i++ {
		port := randomNumberInRange(49152, 61000)
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			l.Close()
			return port, nil
		}
	}

	return 0, errors.New("unable to find open port")
}

// randomNumberInRange returns a random number within the supplied range.
func randomNumberInRange(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
