package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/application"
	"github.com/moneas/bookstore/internal/domain/order"
)

type OrderHandler struct {
	service *application.OrderService
}

func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.service.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order order.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	order.UserID = id

	err := h.service.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderID := order.ID
	orderDetail, _ := h.service.GetOrderDetails(int(orderID))
	c.JSON(http.StatusCreated, orderDetail)
}

func (h *OrderHandler) GetOrderDetails(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.service.GetOrderDetails(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	orders, err := h.service.GetOrdersByUserID(int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orders not found"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User didn't have any orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
