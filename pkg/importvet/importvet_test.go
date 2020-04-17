// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/talos-systems/importvet/pkg/importvet"
)

func TestIntegration(t *testing.T) {
	testdata := analysistest.TestData()
	if err := importvet.InitConfig(testdata); err != nil {
		t.Error(err)
	}

	analysistest.Run(t, testdata, importvet.Analyzer, "github.com/author2/pkg/...")
	analysistest.Run(t, testdata, importvet.Analyzer, "github.com/author3/pkg2/...")
}
