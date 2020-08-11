// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dghubble/trie"
)

// ConfigTree accumulates config files by looking recursively over the tree.
type ConfigTree struct {
	configTrie *trie.PathTrie
}

// NewConfigTree generates config tree by walking the source tree.
func NewConfigTree(rootPath string) (*ConfigTree, error) {
	cfgTree := ConfigTree{}
	cfgTree.configTrie = trie.NewPathTrie()

	var err error

	rootPath, err = filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	if err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() && filepath.Base(path) == configFilename {
			cfg, err := LoadConfig(path)
			if err != nil {
				return err
			}

			cfgTree.configTrie.Put(filepath.Dir(path), cfg)

			return nil
		}

		if info.IsDir() && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &cfgTree, nil
}

// Match returns list of configs applicable to the path.
//
// All matching configs in the tree are returned with least specific
// config first.
func (cfgTree *ConfigTree) Match(path string) ([]*Config, error) {
	var configs []*Config

	err := cfgTree.configTrie.WalkPath(path, func(key string, value interface{}) error {
		configs = append(configs, value.(*Config))

		return nil
	})

	return configs, err
}
