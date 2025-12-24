// SPDX-License-Identifier: GPL-3.0-or-later
package getopt

import (
	"testing"
)

func Test_getopt_long(t *testing.T) {
	tcase := Test_case{
		cfg_optstr: "w:p:hvj:",
		cfg_longopt: []Option{
			{"path",       true,  'w'},
			{"path-to",    true,  'w'},
			{"version",    false, 'v'},
			{"help",       false, 'h'},
			{"help-me",    false, 'h'},
		},
		argv: []string{
			"a.out", "-v", "-h",  "--path", "val1",  "-w", "val2",
			"--path-to", "val3",  "--help-me",
		},
		exps: []Expectation{
			{'v', ""},  {'h', ""},
			{'w', "val1"},  {'w', "val2"},
			{'w', "val3"},  {'h', ""},
		},
	}
	tcase.Test(t);
}

// GNU style options test (e.g. -x666, -XPOST)
func Test_getopt_long_gnu_style(t *testing.T) {
	tcase := Test_case{
		cfg_optstr: "x:hX:",
		cfg_longopt: []Option{
			{"value",      true,  'x'},
			{"method",     true,  'X'},
			{"help",       false, 'h'},
		},
		argv: []string{
			"a.out", "--help",  "-x666", "-XPOST",
		},
		exps: []Expectation{
			{'h', ""},
			{'x', "666"}, {'X', "POST"},
		},
	}
	tcase.Test(t);
}

// Edge cases
func Test_getopt_long_edges(t *testing.T) {
	var tcase Test_case;
	long_opts := []Option{
		{"value",      true,  'x'},
		{"method",     true,  'X'},
		{"help",       false, 'h'},
	};

	tcase = Test_case{ // null long_opts
		cfg_optstr: "x:hX:",
		cfg_longopt: nil,
		exps: []Expectation{},
	}
	tcase.Test(t);
	tcase = Test_case{ // empty optstr
		cfg_optstr: "",
		cfg_longopt: nil,
		exps: []Expectation{},
		}
	tcase.Test(t);
	tcase = Test_case{ // empty input
		cfg_optstr: "x:hX:",
		cfg_longopt: long_opts,
		argv: []string{},
		exps: []Expectation{},
	}
	tcase.Test(t);
	tcase = Test_case{ // no arg
		cfg_optstr: "x:hX:",
		cfg_longopt: long_opts,
		argv: []string{"a.out"},
		exps: []Expectation{},
	}
	tcase.Test(t);
}
	
// Non existing options test
func Test_getopt_nonsexist_opt(t *testing.T) {
	tcase := Test_case{
		cfg_optstr: "x:hX:",
		cfg_longopt:  []Option{
			{"value",      true,  'x'},
			{"method",     true,  'X'},
			{"help",       false, 'h'},
		},
		argv: []string{
			"a.out", "-a",  "-xTEST",  "-bBB", "--non-exist",  "-h",
		},
		exps: []Expectation{
			// -a and -b, --non-exist MUST be ignored
			{'x', "TEST"},
			{'h', ""},
		},
	}
	tcase.Test(t);
}
