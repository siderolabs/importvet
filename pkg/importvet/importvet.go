// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package importvet implements import path linter.
package importvet

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
)

const configFilename = ".importvet.yaml"

// Analyzer ...
var Analyzer = &analysis.Analyzer{
	Name:      "importvet",
	Doc:       "checks that import paths conforman to restrictions",
	Run:       run,
	FactTypes: []analysis.Fact{new(importFact)},
}

var configTree *ConfigTree

// InitConfig should be called to initialize configs for the import restrictions.
func InitConfig(rootPath string) (err error) {
	configTree, err = NewConfigTree(rootPath)

	return err
}

func run(pass *analysis.Pass) (interface{}, error) {
	if configTree == nil {
		return nil, fmt.Errorf("config tree wasn't initialized for importvet")
	}

	// figure out path for the package
	var packagePath string

	for _, f := range pass.Files {
		pos := pass.Fset.Position(f.Package)
		if pos.IsValid() {
			packagePath = filepath.Dir(pos.Filename)

			break
		}
	}

	if packagePath == "" {
		// package path wasn't discovered, skip it
		return nil, nil //nolint:nilnil
	}

	configs, err := configTree.Match(packagePath)
	if err != nil {
		return nil, err
	}

	fact := importFact{Imports: map[*types.Package]struct{}{}}

	verified := map[string]struct{}{}

	for _, f := range pass.Files {
		for _, imp := range f.Imports {
			pkg := imported(pass.TypesInfo, imp)

			for _, config := range configs {
				if config.Process(pkg) == ActionDeny {
					pass.ReportRangef(imp, "import path %v is denied by config %s", imp.Path.Value, config.Path)
				}
			}

			path := pkg.Path()

			verified[path] = struct{}{}

			if _, exists := fact.Imports[pkg]; !exists {
				fact.Imports[pkg] = struct{}{}

				var otherFact importFact
				if len(configs) > 0 && pass.ImportPackageFact(pkg, &otherFact) {
					otherFact.Verify(pass, configs, imp, []string{path}, verified)
				}
			}
		}
	}

	if len(fact.Imports) > 0 {
		pass.ExportPackageFact(&fact)
	}

	return nil, nil //nolint:nilnil
}

func imported(info *types.Info, spec *ast.ImportSpec) *types.Package {
	obj, ok := info.Implicits[spec]
	if !ok {
		obj = info.Defs[spec.Name] // renaming import
	}

	return obj.(*types.PkgName).Imported() //nolint:forcetypeassert
}
