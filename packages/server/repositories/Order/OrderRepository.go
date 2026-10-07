package order

import (
	"path/filepath"

	"gymshark/server/entities"
	base "gymshark/server/repositories/Base"
)

func New(dbDir string) *base.Repository[entities.Order] {
	return base.New[entities.Order](filepath.Join(dbDir, "orders.json"))
}
