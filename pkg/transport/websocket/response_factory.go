package websocket

import (
	"github.com/krls256/dsd2024additional/pkg/entities"
	"github.com/krls256/dsd2024additional/pkg/transport"
	"github.com/samber/mo"
	"net/http"
)

func NewResponseFactory(writeChan chan *transport.Response) *ResponseFactory {
	return &ResponseFactory{writeChan: writeChan}
}

type ResponseFactory struct {
	writeChan chan *transport.Response
}

func (f *ResponseFactory) OK(action string, meta, nonce, data interface{}) {
	f.writeChan <- transport.NewResponse(http.StatusOK, meta, nonce, data)
}

func (f *ResponseFactory) BadRequest(action string, meta, nonce interface{}, err ...error) {
	f.writeChan <- transport.NewFormattedError(http.StatusBadRequest, meta, nonce, err)
}

func (f *ResponseFactory) NotFound(action string, meta, nonce interface{}, err ...error) {
	f.writeChan <- transport.NewFormattedError(http.StatusNotFound, meta, nonce, err)
}

func (f *ResponseFactory) NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		rf:    f,
		nonce: nil,
	}
}

func NewResponseWriter(rf *ResponseFactory, nonce interface{}) *ResponseWriter {
	return &ResponseWriter{
		rf:    rf,
		nonce: nonce,
	}
}

type ResponseWriter struct {
	rf    *ResponseFactory
	nonce interface{}
}

func (w *ResponseWriter) OK(action string, meta, data interface{}) {
	w.rf.OK(action, meta, w.nonce, data)
}

func (w *ResponseWriter) BadRequest(action string, meta interface{}, err ...error) {
	w.rf.BadRequest(action, meta, w.nonce, err...)
}

func (w *ResponseWriter) NotFound(action string, meta interface{}, err ...error) {
	w.rf.NotFound(action, meta, w.nonce, err...)
}

func (w *ResponseWriter) AsyncWriter() *AsyncWriter {
	return &AsyncWriter{
		w: w,
	}
}

func (w *AsyncWriter) AsyncWrite(result mo.Either[entities.LeftWithAction, error]) {
	if l, ok := result.Left(); ok {
		w.w.OK(l.Action, nil, l.Payload)
	}

	if r, ok := result.Right(); ok {
		w.w.BadRequest(InternalErrorAction, nil, r)
	}
}

type AsyncWriter struct {
	w *ResponseWriter
}
