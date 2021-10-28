// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package importvet

import (
	"fmt"
	"go/types"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Action is rule action (firewall style).
type Action string

// Defined actions.
const (
	ActionAllow = Action("allow")
	ActionDeny  = Action("deny")
)

// Rule represents single import restriction rule.
//
//nolint:govet
type Rule struct {
	// Regexp is regular expression against the import path.
	Regexp string `yaml:"regexp"`
	// Set is a way to refer to some set of packages (e.g. `std` packages).
	Set string `yaml:"set"`
	// Action is one of the "Allow", "Deny"
	Action Action `yaml:"action"`
	// Stop indicates whether rule processing should stop or continue.
	Stop bool `yaml:"stop"`

	compiled *regexp.Regexp
	set      packageSet
}

// Matches checks whether rule applies to imported package.
func (rule *Rule) Matches(pkg *types.Package) bool {
	if rule.set != nil {
		if !rule.set.Matches(pkg) {
			return false
		}
	}

	if rule.compiled != nil {
		if !rule.compiled.MatchString(pkg.Path()) {
			return false
		}
	}

	return true
}

// Validate rule syntax.
func (rule *Rule) Validate() error {
	if rule.Regexp == "" && rule.Set == "" {
		return fmt.Errorf("rule regexp and set are empty")
	}

	var err error

	switch {
	case rule.Regexp != "":
		rule.compiled, err = regexp.Compile(rule.Regexp)
		if err != nil {
			return fmt.Errorf("error in regexp syntax %q: %w", rule.Regexp, err)
		}
	case rule.Set != "":
		rule.set, err = lookupPackageSet(rule.Set)
		if err != nil {
			return err
		}
	}

	switch rule.Action {
	case ActionAllow:
	case ActionDeny:
	default:
		return fmt.Errorf("unsupported action %q", rule.Action)
	}

	return nil
}

// RuleSet is an order set of Rules processed top-down.
//
// Rules are processed firewall-style, top-down, only matching rules are applied.
// Action of the last matched rule (or rule with Stop enabled) is the action taken.
// Default action is 'allow'.
type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

// Process evaluates rules and returns final result: is import of pkg allowed or not.
func (ruleset *RuleSet) Process(pkg *types.Package) Action {
	action := ActionAllow

	for _, rule := range ruleset.Rules {
		if rule.Matches(pkg) {
			action = rule.Action

			if rule.Stop {
				break
			}
		}
	}

	return action
}

// Validate the rules, stops on first error.
func (ruleset *RuleSet) Validate() error {
	for i := range ruleset.Rules {
		if err := ruleset.Rules[i].Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Config is YAML representation of the config.
type Config struct {
	Path    string `yaml:"-"`
	RuleSet `yaml:",inline"`
}

// LoadConfig loads import restrictions config from specified path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close() //nolint:errcheck

	cfg := Config{
		Path: path,
	}

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error processing config file %q: %w", path, err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, f.Close()
}
