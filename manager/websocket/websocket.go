package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/observiq/observiq-collector/manager/message"
	"golang.org/x/sync/errgroup"
)

// Conn is a websocket connection.
type Conn = websocket.Conn

// ErrConnectionClose is the error returned during a normal connection closure.
var ErrConnectionClosed = errors.New("connection closed")

// Open will open a new websocket connection.
func Open(ctx context.Context, url string, headers http.Header) (*Conn, error) {
	conn, res, err := websocket.DefaultDialer.DialContext(ctx, url, headers)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	return conn, err
}

// Close will attempt to gracefully close a websocket connection.
func Close(conn *Conn) {
	deadline := time.Now().Add(time.Second)
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")

	err := conn.WriteControl(websocket.CloseMessage, message, deadline)
	if err != nil {
		conn.Close()
	}
}

// HandleReceive will handle receiving inbound traffic from the connection
// until the supplied context is cancelled or an error occurs.
func HandleReceive(ctx context.Context, conn *Conn, in chan *message.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := read(conn)
			if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return fmt.Errorf("unknown read error: %w", err)
			}

			if err != nil {
				return ErrConnectionClosed
			}

			in <- message
		}
	}
}

// HandleSend will handle sending outbound traffic to the connection until
// the supplied context is cancelled or the outbound channel is closed.
func HandleSend(ctx context.Context, conn *Conn, out chan *message.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message, ok := <-out:
			if !ok {
				return nil
			}

			if err := write(message, conn); err != nil {
				return fmt.Errorf("unknown write error: %w", err)
			}
		}
	}
}

// HandleClose will handle closing the connection when context is cancelled.
func HandleClose(ctx context.Context, conn *Conn) error {
	defer Close(conn)
	<-ctx.Done()
	return ctx.Err()
}

// HandleTraffic will handle sending and receiving traffic from the supplied connection.
// If an error occurs or context is cancelled, the connection will be closed.
func HandleTraffic(ctx context.Context, conn *Conn, in, out chan *message.Message) error {
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return HandleReceive(groupCtx, conn, in) })
	group.Go(func() error { return HandleSend(groupCtx, conn, out) })
	group.Go(func() error { return HandleClose(groupCtx, conn) })
	return group.Wait()
}

// read will read the next message from the supplied connection.
func read(conn *Conn) (*message.Message, error) {
	_, reader, err := conn.NextReader()
	if err != nil {
		return nil, err
	}

	message := &message.Message{}
	err = json.NewDecoder(reader).Decode(message)
	return message, err
}

// write will write a message to the supplied connection.
func write(m *message.Message, conn *Conn) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, bytes)
}
