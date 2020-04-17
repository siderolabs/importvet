// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type packageSet map[string]struct{}

func (set packageSet) Matches(pkg *types.Package) bool {
	_, ok := set[pkg.Path()]

	return ok
}

var cachedSets map[string]packageSet

// lookupPackageSet finds package set by name.
func lookupPackageSet(name string) (packageSet, error) {
	if cachedSets == nil {
		cachedSets = map[string]packageSet{}
	}

	set, ok := cachedSets[name]
	if ok {
		return set, nil
	}

	set = packageSet{}

	switch name {
	case "std":
		pkgs, err := packages.Load(nil, "std")
		if err != nil {
			return nil, err
		}

		for _, p := range pkgs {
			set[p.PkgPath] = struct{}{}
		}
	default:
		return nil, fmt.Errorf("unsupported package set %q", name)
	}

	cachedSets[name] = set

	return set, nil
}
