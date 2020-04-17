// want package:`importFact\(\[fmt]\)`
package b

import (
	_ "fmt" // want `import path "fmt" is denied by config`
)
