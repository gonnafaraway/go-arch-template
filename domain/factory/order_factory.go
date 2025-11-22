package factory

import (
	"go-arch-template/domain/entity"
	"go-arch-template/ports/repository"
)

type OrderFactory struct {
	productRepo repository.ProductRepository
}

func NewOrderFactory(productRepo repository.ProductRepository) *OrderFactory {
	return &OrderFactory{
		productRepo: productRepo,
	}
}

func (f *OrderFactory) CreateOrderFromCart(ctx context.Context, userID string, cartItems []CartItem) (*entity.Order, error) {
	productIDs := make([]string, len(cartItems))
	for i, item := range cartItems {
		productIDs[i] = item.ProductID
	}

	products, err := f.productRepo.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	orderItems := make([]entity.OrderItem, len(cartItems))
	for i, cartItem := range cartItems {
		product, exists := products[cartItem.ProductID]
		if !exists {
			return nil, entity.ErrProductNotFound
		}

		orderItems[i] = entity.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
		}
	}

	return entity.NewOrder(userID, orderItems)
}
