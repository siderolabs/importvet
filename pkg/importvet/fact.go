// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type importFact struct {
	Imports map[*types.Package]struct{}
}

func (fact *importFact) AFact() {}

func (fact *importFact) Verify(pass *analysis.Pass, config *Config, rng analysis.Range, chain []string, verified map[string]struct{}) {
	for pkg := range fact.Imports {
		path := pkg.Path()

		if _, ok := verified[path]; ok {
			continue
		}

		verified[path] = struct{}{}

		if config.Process(pkg) == ActionDeny {
			pass.ReportRangef(rng, "import path %v is denied by config (via chain %s)", path, strings.Join(chain, " -> "))
		}

		var innerFact importFact

		if pass.ImportPackageFact(pkg, &innerFact) {
			innerFact.Verify(pass, config, rng, append(chain, pkg.Path()), verified)
		}
	}
}

func (fact *importFact) String() string {
	imports := make([]string, 0, len(fact.Imports))

	for pkg := range fact.Imports {
		imports = append(imports, pkg.Path())
	}

	sort.Strings(imports)

	return fmt.Sprintf("importFact(%v)", imports)
}
