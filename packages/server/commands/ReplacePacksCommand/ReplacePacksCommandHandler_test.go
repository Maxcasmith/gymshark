package replacepackscommand

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gymshark/server/entities"
	packrepo "gymshark/server/repositories/Pack"
)

func TestReplacePacksCommandHandler(t *testing.T) {
	dbDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dbDir, "packs.json"), []byte(`[{"id":1,"size":250}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	packs := packrepo.New(dbDir)
	handler := &ReplacePacksCommandHandler{Packs: packs}

	saved, err := handler.Handle(ReplacePacksCommand{Sizes: []int{53, 23, 31, 23}})
	if err != nil {
		t.Fatal(err)
	}
	expected := []entities.Pack{{ID: 1, Size: 23}, {ID: 2, Size: 31}, {ID: 3, Size: 53}}
	if !reflect.DeepEqual(saved, expected) {
		t.Errorf("saved %+v, expected %+v", saved, expected)
	}
	if stored, _ := packs.All(); !reflect.DeepEqual(stored, expected) {
		t.Errorf("stored %+v, expected %+v", stored, expected)
	}

	tooMany := make([]int, MaxPackSizes+1)
	for index := range tooMany {
		tooMany[index] = index + 1
	}
	for _, sizes := range [][]int{nil, {0}, {-5, 250}, {1_000_001}, tooMany} {
		if _, err := handler.Handle(ReplacePacksCommand{Sizes: sizes}); !errors.Is(err, ErrInvalidPacks) {
			t.Errorf("sizes %v: expected ErrInvalidPacks, got %v", sizes, err)
		}
	}
	if stored, _ := packs.All(); !reflect.DeepEqual(stored, expected) {
		t.Errorf("rejected input changed storage: %+v", stored)
	}

	if _, err := handler.Handle(ReplacePacksCommand{Sizes: DefaultPackSizes}); err != nil {
		t.Errorf("default sizes rejected: %v", err)
	}
}
