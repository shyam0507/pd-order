package api

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shyam0507/pd-order/src/src/storage"
)

type Server struct {
	port     string
	storage  storage.Storage
	producer storage.Producer
	r        *echo.Echo
}

func NewServer(port string, storage storage.Storage, producer storage.Producer) *Server {
	return &Server{port: port, storage: storage, producer: producer, r: echo.New()}
}

func (s *Server) Start() {
	e := s.r
	e.Use(middleware.RequestID())

	g := e.Group("/api/orders/v1.0")

	g.POST("/", s.createOrder)

	slog.Info("Server started on port", "Port", s.port)

	e.Start(":" + s.port)

}
