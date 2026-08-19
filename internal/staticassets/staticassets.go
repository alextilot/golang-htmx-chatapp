// Package staticassets builds cache-busting URLs for files served under
// /static, so a long-lived, immutable Cache-Control header (see
// internal/router) can be paired with a URL that changes whenever a file's
// content changes.
package staticassets

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"
)

const root = "web/static"

var (
	mu    sync.Mutex
	cache = map[string]string{}
)

// URL returns the /static URL for relPath (relative to web/static), with a
// content-hash query parameter appended.
func URL(relPath string) string {
	return "/static/" + relPath + "?v=" + version(relPath)
}

func version(relPath string) string {
	mu.Lock()
	defer mu.Unlock()

	if v, ok := cache[relPath]; ok {
		return v
	}

	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		return "dev"
	}

	sum := sha256.Sum256(data)
	v := hex.EncodeToString(sum[:])[:8]
	cache[relPath] = v
	return v
}
