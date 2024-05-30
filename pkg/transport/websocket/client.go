package websocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/krls256/dsd2024additional/pkg/transport"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var (
	ErrContextDeadlineExceeded = errors.New("context deadline exceeded")
	ErrUnknownAction           = errors.New("unknown action")
	ErrClosedPipe              = errors.New("closed pipe")
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	defaultWriteBufferSize = 100
	writeWait              = 10 * time.Second
	pongWait               = 60 * time.Second
	pingPeriod             = (pongWait * 9) / 10
	maxMessageSize         = 512
)

func NewClient(conn *websocket.Conn, router *Router,
	wg *sync.WaitGroup) *Client {
	wc := make(chan *transport.Response, defaultWriteBufferSize)

	return &Client{
		conn:   conn,
		router: router,

		writeChan: wc,
		rf:        NewResponseFactory(wc),

		wg:    wg,
		close: make(chan struct{}),
	}
}

type Client struct {
	conn   *websocket.Conn
	router *Router

	writeChan chan *transport.Response

	rf *ResponseFactory

	wg        *sync.WaitGroup
	close     chan struct{}
	closeOnce sync.Once
}

func (c *Client) Stop() {
	c.closeOnce.Do(func() {
		close(c.close)
	})
}

func (c *Client) reader(idleOnlineWait time.Duration) {
	defer c.wg.Done()
	defer c.Stop()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(c.pongHandler)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		zap.S().Debug(err)
	}

	timer := time.NewTimer(idleOnlineWait)
	defer timer.Stop()

	for {
		select {
		case message, ok := <-c.readMessageChan():
			if !ok {
				return
			}

			message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

			c.read(message)
		case <-c.close:
			return
		case <-timer.C:
			return
		}

		timer.Reset(idleOnlineWait)
	}
}

func (c *Client) read(bts []byte) {
	req := &Request{}

	if err := json.Unmarshal(bts, req); err != nil {
		zap.S().Error(err)
		c.rf.BadRequest(InternalErrorAction, nil, req.Nonce, err)

		return
	}

	go c.handle(req)
}

func (c *Client) writer() {
	ticker := time.NewTicker(pingPeriod)

	defer c.wg.Done()
	defer c.Stop()
	defer ticker.Stop()
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.writeChan:
			if err := c.writeHandler(message, ok); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.tickHandler(); err != nil {
				return
			}
		case <-c.close:
			return
		}
	}
}

func (c *Client) write(event interface{}) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	bts, err := json.Marshal(event)
	if err != nil {
		return
	}

	if _, err = w.Write(bts); err != nil {
		zap.S().Error(err)
	}

	if err = w.Close(); err != nil {
		zap.S().Error("can't write to closed websocket conn", err)

		return
	}
}

func (c *Client) writeHandler(message *transport.Response, ok bool) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		zap.S().Debug(err)
	}

	if !ok {
		if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			zap.S().Debug(err)
		}

		return ErrClosedPipe
	}

	c.write(message)

	return nil
}

func (c *Client) readMessageChan() chan []byte {
	ch := make(chan []byte)

	go func() {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			zap.S().Error(err)
			close(ch)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Debug(err)
			}

			return
		}
		ch <- message
	}()

	return ch
}

func (c *Client) pongHandler(string) error {
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		zap.S().Error(err)
	}

	return nil
}

func (c *Client) tickHandler() error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		zap.S().Debug(err)
	}

	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) handle(req *Request) {
	done := make(chan struct{})
	ctx := c.newContext(req)

	go func() {
		handleFunc, ok := c.router.find(req.Action)
		if !ok {
			c.rf.NotFound(req.Action, nil, req.Nonce, ErrUnknownAction)

			return
		}

		handleFunc(ctx)

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		c.rf.BadRequest(req.Action, req.Nonce, ErrContextDeadlineExceeded)

	case <-done:
	}
}

func (c *Client) newContext(req *Request) *Context {
	return &Context{
		Context: context.Background(),
		req:     req,
		rf:      c.rf,
	}
}
