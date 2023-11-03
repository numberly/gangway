package templates_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/numberly/gangway/templates"
)

func TestTemplateFS(t *testing.T) {
	t.Run("finds templates", func(t *testing.T) {
		filenames := []string{"commandline.tmpl", "home.tmpl"}
		var missing, empty []string

		for _, filename := range filenames {
			file, err := templates.FS.ReadFile(filename)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					missing = append(missing, filename)
					continue
				}
				t.Fatalf("failed to open template %s: %v", filename, err)
			}
			if len(file) == 0 {
				empty = append(empty, filename)
			}
		}
		if len(missing) > 0 {
			t.Fatalf("couldn't find templates: %v", missing)
		}
		if len(empty) > 0 {
			t.Fatalf("empty templates: %v", empty)
		}
	})

	t.Run("does not include non-templates", func(t *testing.T) {
		filenames := []string{"fs.go", "fs_test.go"}
		var found []string

		for _, filename := range filenames {
			_, err := templates.FS.ReadFile(filename)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				t.Fatalf("failed to open template %s: %v", filename, err)
			}
			found = append(found, filename)
		}
		if len(found) > 0 {
			t.Fatalf("found files: %v", found)
		}
	})
}
