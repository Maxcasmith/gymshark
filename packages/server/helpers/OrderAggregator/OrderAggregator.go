package orderaggregator

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"gymshark/server/entities"
)

const MaxRequestedItems = 1_000_000

var ErrInvalidOrder = errors.New("invalid order")

func OrderAggregator(requestedItems int, packSizes []int) (entities.Order, error) {
	if requestedItems < 1 || requestedItems > MaxRequestedItems {
		return entities.Order{}, fmt.Errorf("%w: items must be between 1 and %d", ErrInvalidOrder, MaxRequestedItems)
	}
	if len(packSizes) == 0 {
		return entities.Order{}, errors.New("no pack sizes configured")
	}
	smallestPackSize := slices.Min(packSizes)
	if smallestPackSize < 1 {
		return entities.Order{}, fmt.Errorf("pack sizes must be positive, got %d", smallestPackSize)
	}

	searchLimit := requestedItems + smallestPackSize
	fewestPacks := make([]int, searchLimit)
	lastPackSize := make([]int, searchLimit)
	for total := 1; total < searchLimit; total++ {
		for _, packSize := range packSizes {
			previousTotal := total - packSize
			if previousTotal < 0 || (previousTotal > 0 && fewestPacks[previousTotal] == 0) {
				continue
			}
			if fewestPacks[total] == 0 || fewestPacks[previousTotal]+1 < fewestPacks[total] {
				fewestPacks[total] = fewestPacks[previousTotal] + 1
				lastPackSize[total] = packSize
			}
		}
	}

	shippedItems := requestedItems
	for fewestPacks[shippedItems] == 0 {
		shippedItems++
	}

	quantityBySize := map[int]int{}
	for total := shippedItems; total > 0; total -= lastPackSize[total] {
		quantityBySize[lastPackSize[total]]++
	}
	packs := make([]entities.OrderPack, 0, len(quantityBySize))
	for packSize, quantity := range quantityBySize {
		packs = append(packs, entities.OrderPack{PackSize: packSize, Quantity: quantity})
	}
	slices.SortFunc(packs, func(left, right entities.OrderPack) int {
		return cmp.Compare(right.PackSize, left.PackSize)
	})

	return entities.Order{
		RequestedItems: requestedItems,
		ShippedItems:   shippedItems,
		ExcessItems:    shippedItems - requestedItems,
		TotalPacks:     fewestPacks[shippedItems],
		Packs:          packs,
	}, nil
}
