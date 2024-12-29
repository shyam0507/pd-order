package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shyam0507/pd-order/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) createOrder(c echo.Context) error {
	slog.Info("Creating order")
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

	event := types.OrderCreatedEvent{
		Id:              order.Id.Hex(),
		Type:            "OrderCreated",
		Source:          "OrderService",
		SpecVersion:     "1.0",
		Data:            order,
		DataContentType: "application/json",
	}

	s.producer.ProduceOrderCreated(order.Id.Hex(), event)

	return c.JSON(http.StatusOK, map[string]string{
		"orderId": order.Id.Hex(),
	})
}
