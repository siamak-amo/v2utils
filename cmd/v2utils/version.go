// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"

	core "github.com/xtls/xray-core/core"
)

// V2utils current version
const Version = "1.5";

// Prints version of v2ray and xray
func printVersion() {
	fmt.Fprintf (os.Stderr,`V2utils %s - Xray-Core compatible utility.

Compiled with:
`, Version);
	version := core.VersionStatement();
	for _, s := range version {
		fmt.Fprintf (os.Stderr, "    %s\n", s);
	}
}
