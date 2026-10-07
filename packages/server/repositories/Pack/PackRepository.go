package pack

import (
	"path/filepath"

	"gymshark/server/entities"
	base "gymshark/server/repositories/Base"
)

func New(dbDir string) *base.Repository[entities.Pack] {
	return base.New[entities.Pack](filepath.Join(dbDir, "packs.json"))
}
