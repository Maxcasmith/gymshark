package entities

type OrderPack struct {
	PackSize int `json:"packSize"`
	Quantity int `json:"quantity"`
}

type Order struct {
	ID             int         `json:"id"`
	RequestedItems int         `json:"requestedItems"`
	ShippedItems   int         `json:"shippedItems"`
	ExcessItems    int         `json:"excessItems"`
	TotalPacks     int         `json:"totalPacks"`
	Packs          []OrderPack `json:"packs"`
}
