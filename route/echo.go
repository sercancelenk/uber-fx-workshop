package route

import (
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"uber-fx-workshop/service"
)

type EchoHandler struct {
	log *zap.Logger
	c1  *service.Consumer1
}

func NewEchoHandler(log *zap.Logger, c1 *service.Consumer1) *EchoHandler {
	return &EchoHandler{
		log: log,
		c1:  c1,
	}
}
func (e *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(8 * time.Second)
	_ = e.c1.Consume("message received 1")
	if _, err := io.Copy(w, r.Body); err != nil {
		e.log.Warn("Failed to handle request", zap.Error(err))
	}
}

func (e *EchoHandler) Pattern() string {
	return "/echo"
}
