package order

type Repository interface {
	FindAll() ([]Order, error)
	FindByID(id uint) (*Order, error)
	Save(order *Order) error
	FindOrderDetails(orderID int) (*OrderDetail, error)
	FindOrdersByUserID(userID int) ([]*OrderDetail, error)
}
