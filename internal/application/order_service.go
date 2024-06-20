package application

import (
	"github.com/moneas/bookstore/internal/domain/order"
)

type OrderService struct {
	repo order.Repository
}

func NewOrderService(repo order.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetOrders() ([]order.Order, error) {
	return s.repo.FindAll()
}

func (s *OrderService) GetOrderByID(id uint) (*order.Order, error) {
	return s.repo.FindByID(id)
}

func (s *OrderService) CreateOrder(order *order.Order) error {
	return s.repo.Save(order)
}

func (s *OrderService) GetOrderDetails(orderID int) (*order.OrderDetail, error) {
	return s.repo.FindOrderDetails(orderID)
}

func (s *OrderService) GetOrdersByUserID(userID int) ([]*order.OrderDetail, error) {
	return s.repo.FindOrdersByUserID(userID)
}
