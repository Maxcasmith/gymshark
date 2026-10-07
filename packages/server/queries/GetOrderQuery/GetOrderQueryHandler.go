package getorderquery

import (
	"errors"

	"gymshark/server/entities"
	querybase "gymshark/server/queries/Base"
	repobase "gymshark/server/repositories/Base"
)

var ErrNotFound = errors.New("order not found")

type GetOrderQueryHandler struct {
	Orders repobase.RepositoryInterface[entities.Order]
}

var _ querybase.QueryHandlerInterface[GetOrderQuery, entities.Order] = (*GetOrderQueryHandler)(nil)

func (handler *GetOrderQueryHandler) Handle(query GetOrderQuery) (entities.Order, error) {
	orders, err := handler.Orders.All()
	if err != nil {
		return entities.Order{}, err
	}
	for _, order := range orders {
		if order.ID == query.OrderID {
			return order, nil
		}
	}
	return entities.Order{}, ErrNotFound
}
