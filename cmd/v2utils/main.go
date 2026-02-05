// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
)

func main() {
	opt := Opt{};
	opt.GetArgs();
	if ret := opt.HandleArgs(); ret < 0 {
		os.Exit (-ret);
	}
	if ret := opt.Init(); ret < 0 {
		os.Exit (-ret);
	}

	// The main loop
	for ;; {
		if EOF := opt.GetInput(); true == EOF {
			break;
		}
		if opt.Do() < 0 {
			break;
		}
	}
}
