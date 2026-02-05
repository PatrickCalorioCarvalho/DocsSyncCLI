package scanner

import (
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/config"
	"os"
)

func Scan(cfg *config.Config, projectRoot string) ([]string, error) {
	root := filepath.Join(projectRoot, cfg.Scan.Root)

	var results []string

	for _, pattern := range cfg.Scan.Include {
		matches, err := doublestar.Glob(os.DirFS(root), pattern)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			fullPath := filepath.Join(root, match)

			if shouldExclude(match, cfg.Scan.Exclude) {
				continue
			}

			results = append(results, fullPath)
		}
	}

	return results, nil
}

func shouldExclude(path string, excludes []string) bool {
	for _, pattern := range excludes {
		ok, _ := doublestar.Match(pattern, path)
		if ok {
			return true
		}
	}
	return false
}
