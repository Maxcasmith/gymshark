package createordercommand

type CreateOrderCommand struct {
	RequestedItems int `json:"requestedItems"`
}

func (CreateOrderCommand) CommandName() string { return "CreateOrderCommand" }
