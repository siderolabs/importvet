// want package:`importFact\(\[fmt github.com/author2/pkg/dir1 github.com/author2/pkg/dir2 github.com/author3/pkg2 github.com/author4/pkg3\]\)`
package main

import (
	_ "fmt" // want `import path "fmt" is denied by config`

	_ "github.com/author2/pkg/dir1"
	_ "github.com/author2/pkg/dir2"
	_ "github.com/author3/pkg2" // want `import path "github.com/author3/pkg2" is denied by config`
	_ "github.com/author4/pkg3" //want `import path net is denied by config \(via chain github.com/author4/pkg3\)`
)
