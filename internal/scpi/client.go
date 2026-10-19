// Package scpi provides communication with SCPI instruments over TCP.
package scpi

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client owns one TCP session. Calls must be sequential. A binary query is the
// final query on a session: its optional trailing terminator is left unread.
// Close the client after QueryBlock, including on error.
type Client struct {
	conn         net.Conn
	cleanup      func()
	blockQueried bool
}

// Dial bounds the whole session (connect, identification and screenshot) by timeout.
func Dial(ctx context.Context, address string, timeout time.Duration) (*Client, error) {
	conn, cleanup, err := connect(ctx, address, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, cleanup: cleanup}, nil
}

func (c *Client) Close() { c.cleanup() }

// Identify opens a short session for callers that only need raw identification.
func Identify(ctx context.Context, address string, timeout time.Duration) (string, error) {
	client, err := Dial(ctx, address, timeout)
	if err != nil {
		return "", err
	}
	defer client.Close()
	return client.Identify()
}

func (c *Client) Identify() (string, error) {
	if c.blockQueried {
		return "", fmt.Errorf("binary query completed: open a new session")
	}
	if _, err := fmt.Fprint(c.conn, "*IDN?\n"); err != nil {
		return "", fmt.Errorf("send *IDN?: %w", err)
	}
	line, err := bufio.NewReaderSize(c.conn, 4096).ReadSlice('\n')
	if err != nil {
		return "", fmt.Errorf("read *IDN? response (maximum 4096 bytes): %w", err)
	}
	id := strings.TrimSpace(string(line))
	if id == "" {
		return "", fmt.Errorf("empty *IDN? response")
	}
	return id, nil
}

// QueryBlock is transport-only; device profiles own the commands and size limits.
func (c *Client) QueryBlock(command string, maxBytes int) ([]byte, error) {
	if c.blockQueried {
		return nil, fmt.Errorf("only one binary query is permitted per session")
	}
	if command == "" || strings.ContainsAny(command, "\r\n") || maxBytes <= 0 {
		return nil, fmt.Errorf("invalid block query or size limit")
	}
	c.blockQueried = true
	if _, err := fmt.Fprintf(c.conn, "%s\n", command); err != nil {
		return nil, fmt.Errorf("send %s: %w", command, err)
	}
	return ReadBlock(c.conn, maxBytes)
}

// Session deadlines and cancellation also cover blocking reads and writes.
func connect(ctx context.Context, address string, timeout time.Duration) (net.Conn, func(), error) {
	if timeout <= 0 {
		return nil, nil, fmt.Errorf("timeout must be positive")
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("connect to %s: %w", address, err)
	}
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	cleanup := func() { stop(); cancel(); conn.Close() }
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("set deadline: %w", err)
	}
	return conn, cleanup, nil
}
