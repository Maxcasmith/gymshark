package replacepackscommand

import (
	"errors"
	"fmt"
	"slices"

	commandbase "gymshark/server/commands/Base"
	"gymshark/server/entities"
	orderaggregator "gymshark/server/helpers/OrderAggregator"
	repobase "gymshark/server/repositories/Base"
)

const MaxPackSizes = 20

var DefaultPackSizes = []int{250, 500, 1000, 2000, 5000}

var ErrInvalidPacks = errors.New("invalid pack sizes")

type ReplacePacksCommandHandler struct {
	Packs repobase.RepositoryInterface[entities.Pack]
}

var _ commandbase.CommandHandlerInterface[ReplacePacksCommand, []entities.Pack] = (*ReplacePacksCommandHandler)(nil)

func (handler *ReplacePacksCommandHandler) Handle(command ReplacePacksCommand) ([]entities.Pack, error) {
	sizes := slices.Compact(slices.Sorted(slices.Values(command.Sizes)))
	if len(sizes) == 0 || len(sizes) > MaxPackSizes {
		return nil, fmt.Errorf("%w: between 1 and %d sizes required", ErrInvalidPacks, MaxPackSizes)
	}
	if sizes[0] < 1 || sizes[len(sizes)-1] > orderaggregator.MaxRequestedItems {
		return nil, fmt.Errorf("%w: each size must be between 1 and %d", ErrInvalidPacks, orderaggregator.MaxRequestedItems)
	}

	packs := make([]entities.Pack, len(sizes))
	for index, size := range sizes {
		packs[index] = entities.Pack{ID: index + 1, Size: size}
	}
	return packs, handler.Packs.Replace(packs)
}
