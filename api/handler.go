package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shyam0507/pd-order/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) createOrder(c echo.Context) error {
	var order types.Order
	if err := c.Bind(&order); err != nil {
		slog.Error("Error while creating the order", "Err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	order.Status = "CREATED"
	order.Id = primitive.NewObjectID()

	if err := s.storage.CreateOrder(order); err != nil {
		slog.Error("Error while creating the order", "Err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	slog.Info("Order created successfully")

	s.producer.ProduceOrderCreated(order.Id.Hex(), order)

	return c.JSON(http.StatusOK, map[string]string{})
}
