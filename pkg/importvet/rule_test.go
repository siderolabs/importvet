// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet_test

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/talos-systems/importvet/pkg/importvet"
)

type RulesSuite struct {
	suite.Suite
}

func (suite *RulesSuite) TestLoadGoodConfig() {
	cfg, err := importvet.LoadConfig("testdata/rules/good.yaml")
	suite.Require().NoError(err)

	suite.Assert().Equal(importvet.ActionAllow, cfg.Process(types.NewPackage("github.com/example/a", "a")))
	suite.Assert().Equal(importvet.ActionDeny, cfg.Process(types.NewPackage("gopkg.in/yaml.v3", "yaml")))
}

func (suite *RulesSuite) TestLoadBadConfig1() {
	_, err := importvet.LoadConfig("testdata/rules/bad1.yaml")
	suite.Require().Error(err)

	suite.Assert().EqualError(err, "rule regexp and set are empty")
}

func (suite *RulesSuite) TestLoadBadConfig2() {
	_, err := importvet.LoadConfig("testdata/rules/bad2.yaml")
	suite.Require().Error(err)

	suite.Assert().EqualError(err, "error in regexp syntax \"^$[ab\": error parsing regexp: missing closing ]: `[ab`")
}

func (suite *RulesSuite) TestLoadBadConfig3() {
	_, err := importvet.LoadConfig("testdata/rules/bad3.yaml")
	suite.Require().Error(err)

	suite.Assert().EqualError(err, "unsupported action \"flip\"")
}

func (suite *RulesSuite) TestLoadBadConfig4() {
	_, err := importvet.LoadConfig("testdata/rules/bad4.yaml")
	suite.Require().Error(err)

	suite.Assert().EqualError(err, "unsupported package set \"nonstd\"")
}

func TestRulesSuite(t *testing.T) {
	suite.Run(t, new(RulesSuite))
}
