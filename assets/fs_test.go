package assets_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/numberly/gangway/assets"
)

func TestAssetFS(t *testing.T) {
	t.Run("finds gangway assets", func(t *testing.T) {
		filenames := []string{"gangway.css", "gangway.js"}
		var missing, empty []string

		for _, filename := range filenames {
			file, err := assets.FS.ReadFile(filename)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					missing = append(missing, filename)
					continue
				}
				t.Fatalf("failed to open asset %s: %v", filename, err)
			}
			if len(file) == 0 {
				empty = append(empty, filename)
			}
		}
		if len(missing) > 0 {
			t.Fatalf("couldn't find asset: %v", missing)
		}
		if len(empty) > 0 {
			t.Fatalf("empty assets: %v", empty)
		}
	})

	t.Run("does not include non-assets", func(t *testing.T) {
		filenames := []string{"fs.go", "fs_test.go"}
		var found []string

		for _, filename := range filenames {
			_, err := assets.FS.ReadFile(filename)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				t.Fatalf("failed to open asset %s: %v", filename, err)
			}
			found = append(found, filename)
		}
		if len(found) > 0 {
			t.Fatalf("found files: %v", found)
		}
	})
}
