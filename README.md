# importvet

`importvet` is a Go linter which verifies that import paths in Go packages conform with the supplied specification.

## Install

Install `importvet`:

    go get github.com/talos-systems/importvet@latest

## Usage

Run `importvet` with Go package pattern to analyze:

    importvet ./...

For each import which violates the specification, `importvet` reports diagnostic information:

    ./talos/internal/pkg/provision/request.go:11:2: import path "github.com/talos-systems/talos/internal/pkg/runtime" is denied by config

`importvet` analyzes not only direct imports in the specified package pattern, but also transitive dependencies, that is package `A` imports package `B` (which is allowed by the specification),
but package `B` imports package `C` which is denied, such import is reported:

    ./talos/internal/pkg/provision/providers/factory.go:11:2: import path github.com/talos-systems/talos/internal/pkg/containers/cri/containerd is denied by config (via chain github.com/talos-systems/talos/internal/pkg/provision -> github.com/talos-systems/talos/pkg/config/types/v1alpha1/generate -> github.com/talos-systems/talos/pkg/config/types/v1alpha1)

## Rules

Import specification is configured using rules which act like firewall: rules are executed top-down,
last matching rule action is used. Rules evaluation might be stopped with `stop: true`.

Rules should be placed into file `.importvet.yaml` at any project subdirectory. `importvet` scans
current directory for rules before linting is started, and applies coonfiguration to matching
packages by looking up closest rule file up in the tree. If more one rule file is found,
`importvet` reports that as error.

Rule format:

```yaml
action: allow|deny # action is required
stop: true # default is false
regexp: ^github.com/some/path # at least one of 'regexp', 'set' is required
set: std
```

Property `regexp` matches against canonical package import path. Property `set` allows to pull in
some predefined package sets, at the moment the only set implemented is `std` which expands to list
Go standard library packages.

Rule matches import path if both `set` and `regexp` match the import path (but empty `set` or
`regexp` matches any import path).

Unless `stop: true` is specified, rule evaluation continues to the end of the rule list.

Default action is `allow` (if no rule matches or `.importvet.yaml` is missing).

Full `.importvet.yaml` looks like:

```yaml
# .importvet.yaml
rules:
  - action: deny
    regexp: .+  - action: allow
    regexp: ^github.com/some/package/.+
  - action: allow
    set: std
```

## golangci-lint

`imporvet` can be used as [golangci-lint](https://github.com/golangci/golangci-lint) plugin.

Details TBD.
