package createordercommand

import (
	commandbase "gymshark/server/commands/Base"
	"gymshark/server/entities"
	orderaggregator "gymshark/server/helpers/OrderAggregator"
	repobase "gymshark/server/repositories/Base"
)

type CreateOrderCommandHandler struct {
	Packs  repobase.RepositoryInterface[entities.Pack]
	Orders repobase.RepositoryInterface[entities.Order]
}

var _ commandbase.CommandHandlerInterface[CreateOrderCommand, entities.Order] = (*CreateOrderCommandHandler)(nil)

func (handler *CreateOrderCommandHandler) Handle(command CreateOrderCommand) (entities.Order, error) {
	packs, err := handler.Packs.All()
	if err != nil {
		return entities.Order{}, err
	}

	packSizes := make([]int, len(packs))
	for index, pack := range packs {
		packSizes[index] = pack.Size
	}

	order, err := orderaggregator.OrderAggregator(command.RequestedItems, packSizes)
	if err != nil {
		return entities.Order{}, err
	}

	return handler.Orders.Append(func(existing []entities.Order) entities.Order {
		order.ID = len(existing) + 1
		return order
	})
}
