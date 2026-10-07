package base

import "gymshark/server/entities"

type RepositoryInterface[TEntity entities.Entity] interface {
	All() ([]TEntity, error)
	Append(build func(existing []TEntity) TEntity) (TEntity, error)
	Replace(entities []TEntity) error
}
