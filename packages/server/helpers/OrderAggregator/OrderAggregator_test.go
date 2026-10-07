package orderaggregator

import (
	"errors"
	"reflect"
	"testing"

	"gymshark/server/entities"
)

func pack(packSize, quantity int) entities.OrderPack {
	return entities.OrderPack{PackSize: packSize, Quantity: quantity}
}

func TestOrderAggregator(t *testing.T) {
	standardSizes := []int{250, 500, 1000, 2000, 5000}
	testCases := []struct {
		requestedItems int
		packSizes      []int
		expectedPacks  []entities.OrderPack
	}{
		{1, standardSizes, []entities.OrderPack{pack(250, 1)}},
		{250, standardSizes, []entities.OrderPack{pack(250, 1)}},
		{251, standardSizes, []entities.OrderPack{pack(500, 1)}},
		{501, standardSizes, []entities.OrderPack{pack(500, 1), pack(250, 1)}},
		{12001, standardSizes, []entities.OrderPack{pack(5000, 2), pack(2000, 1), pack(250, 1)}},
		{500000, []int{23, 31, 53}, []entities.OrderPack{pack(53, 9429), pack(31, 7), pack(23, 2)}},
	}
	for _, testCase := range testCases {
		order, err := OrderAggregator(testCase.requestedItems, testCase.packSizes)
		if err != nil {
			t.Fatalf("requested %d: %v", testCase.requestedItems, err)
		}
		if !reflect.DeepEqual(order.Packs, testCase.expectedPacks) {
			t.Errorf("requested %d: got %+v, expected %+v", testCase.requestedItems, order.Packs, testCase.expectedPacks)
		}
	}

	if _, err := OrderAggregator(0, standardSizes); !errors.Is(err, ErrInvalidOrder) {
		t.Errorf("requested 0: expected ErrInvalidOrder, got %v", err)
	}
}
