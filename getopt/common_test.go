// SPDX-License-Identifier: GPL-3.0-or-later
package getopt

import (
	"fmt"
	"runtime"
	"testing"
)

type Expectation struct {
	Value byte
	Optarg string
}

type Test_case struct {
	cfg_optstr string
	cfg_longopt []Option
	argv []string
	exps []Expectation
};

func fileline() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf ("%s:%d", file, line)
	}
	return ""
}

func (tc *Test_case) Test(t *testing.T) {
	i := 0
	Getopt_reset();
	for idx := 0; idx != -1; i += 1 {
		idx = Getopt_long(tc.argv, tc.cfg_optstr, tc.cfg_longopt)
		if idx == -1 {
			if i < len(tc.exps) {
				t.Fatalf("\n%s:  Test #%d f:ailed - get_opt ended too early.\n", fileline(), i);
			}
			return;
		}
		if i >= len(tc.exps) {
			t.Fatalf("\n%s:  Test #%d failed - not enough tests.\n", fileline(), i);
		}
		if tc.exps[i].Value != (byte)(idx) {
			t.Fatalf(`
%s:  Test #%d failed - wrong get_opt return value:
     (expected: '%c') != (actual: '%c')
`,
				fileline(), i, tc.exps[i].Value, (byte)(idx),
			)
			return
		}
		if tc.exps[i].Optarg != Optarg {
			t.Fatalf(`
%s:  Test #%d failed - wrong optarg value:
     (expected: '%s') != (actual: '%s')
`,
				fileline(), i, tc.exps[i].Optarg, Optarg,
			)
			return
		}
	}
}
