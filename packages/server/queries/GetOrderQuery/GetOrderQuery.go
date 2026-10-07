package getorderquery

type GetOrderQuery struct {
	OrderID int
}

func (GetOrderQuery) QueryName() string { return "GetOrderQuery" }
