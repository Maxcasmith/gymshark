package base

import (
	"encoding/json"
	"os"
	"sync"

	"gymshark/server/entities"
)

type Repository[TEntity entities.Entity] struct {
	filePath string
	mutex    sync.Mutex
}

var _ RepositoryInterface[entities.Pack] = (*Repository[entities.Pack])(nil)

func New[TEntity entities.Entity](filePath string) *Repository[TEntity] {
	return &Repository[TEntity]{filePath: filePath}
}

func (repository *Repository[TEntity]) All() ([]TEntity, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	return repository.read()
}

func (repository *Repository[TEntity]) Append(build func(existing []TEntity) TEntity) (TEntity, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	var empty TEntity
	existing, err := repository.read()
	if err != nil {
		return empty, err
	}

	entity := build(existing)
	return entity, repository.write(append(existing, entity))
}

func (repository *Repository[TEntity]) Replace(entities []TEntity) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	return repository.write(entities)
}

func (repository *Repository[TEntity]) write(entities []TEntity) error {
	contents, err := json.MarshalIndent(entities, "", "  ")
	if err != nil {
		return err
	}

	tempPath := repository.filePath + ".tmp"
	if err := os.WriteFile(tempPath, contents, 0o644); err != nil {
		return err
	}
	return os.Rename(tempPath, repository.filePath)
}

func (repository *Repository[TEntity]) read() ([]TEntity, error) {
	contents, err := os.ReadFile(repository.filePath)
	if err != nil {
		return nil, err
	}
	var stored []TEntity
	if err := json.Unmarshal(contents, &stored); err != nil {
		return nil, err
	}
	return stored, nil
}
