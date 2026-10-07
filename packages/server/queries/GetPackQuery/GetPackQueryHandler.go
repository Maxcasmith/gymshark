package getpackquery

import (
	"gymshark/server/entities"
	querybase "gymshark/server/queries/Base"
	repobase "gymshark/server/repositories/Base"
)

type GetPackQueryHandler struct {
	Packs repobase.RepositoryInterface[entities.Pack]
}

var _ querybase.QueryHandlerInterface[GetPackQuery, []entities.Pack] = (*GetPackQueryHandler)(nil)

func (handler *GetPackQueryHandler) Handle(GetPackQuery) ([]entities.Pack, error) {
	return handler.Packs.All()
}
