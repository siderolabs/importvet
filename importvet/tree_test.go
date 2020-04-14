// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet_test

import (
	"go/types"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/talos-systems/importvet/importvet"
)

type TreeSuite struct {
	suite.Suite
}

func (suite *TreeSuite) TestMatch() {
	cfgTree, err := importvet.NewConfigTree("./testdata/src")
	suite.Require().NoError(err)

	absRoot, err := filepath.Abs("./testdata/src")
	suite.Require().NoError(err)

	matches, err := cfgTree.Match(filepath.Join(absRoot, "some-package"))
	suite.Require().NoError(err)
	suite.Require().Nil(matches)

	matches, err = cfgTree.Match(filepath.Join(absRoot, "github.com/author/package/subdir"))
	suite.Require().NoError(err)
	suite.Require().Len(matches, 2)

	suite.Assert().Equal(importvet.ActionAllow, matches[0].Process(types.NewPackage("github.com/example/pkg", "pkg")))
	suite.Assert().Equal(importvet.ActionDeny, matches[1].Process(types.NewPackage("github.com/example/pkg", "pkg")))

	matches, err = cfgTree.Match(filepath.Join(absRoot, "github.com/author/package"))
	suite.Require().NoError(err)
	suite.Require().Len(matches, 2)

	suite.Assert().Equal(importvet.ActionAllow, matches[0].Process(types.NewPackage("github.com/example/pkg", "pkg")))
	suite.Assert().Equal(importvet.ActionDeny, matches[1].Process(types.NewPackage("github.com/example/pkg", "pkg")))

	matches, err = cfgTree.Match(filepath.Join(absRoot, "github.com/author"))
	suite.Require().NoError(err)
	suite.Require().Len(matches, 1)

	suite.Assert().Equal(importvet.ActionAllow, matches[0].Process(types.NewPackage("github.com/example/pkg", "pkg")))
}

func TestTreeSuite(t *testing.T) {
	suite.Run(t, new(TreeSuite))
}
