// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// package main is the entrypoint for the importvet command.
package main

import (
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/siderolabs/importvet/pkg/importvet"
)

func main() {
	if err := importvet.InitConfig("."); err != nil {
		log.Fatal(err)
	}

	singlechecker.Main(importvet.Analyzer)
}
