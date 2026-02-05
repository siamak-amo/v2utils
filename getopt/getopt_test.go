// SPDX-License-Identifier: GPL-3.0-or-later
package getopt

import (
	"testing"
)

func Test_getopt_long(t *testing.T) {
	Getopt_reset();
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

// long option = value test
func Test_getopt_long_equal(t *testing.T) {
	Getopt_reset();
	tcase := Test_case{
		cfg_optstr: "x:",
		cfg_longopt: []Option{
			{"method", true, 'x'},
		},
		argv: []string{
			"a.out", "--method=POST",  "--method=", "GET",
		},
		exps: []Expectation{
			{'x', "POST"},
			{'x', ""}, // the second --method= is empty
		},
	}
	tcase.Test(t);
}

// GNU style options test (e.g. -x666, -XPOST)
func Test_getopt_long_gnu_style(t *testing.T) {
	Getopt_reset();
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
	Getopt_reset();
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
	
// non-existing and unexpected options test
func Test_getopt_nonsexist_opt(t *testing.T) {
	Getopt_reset();
	tcase := Test_case{
		cfg_optstr: "x:hX:",
		cfg_longopt:  []Option{
			{"value",      true,  'x'},
			{"method",     true,  'X'},
			{"help",       false, 'h'},
		},
		argv: []string{
			"a.out",
			"XXX", // unexpected
			"-a",  // non-existing
			"-xTEST",
			"-b", "XXX", "--non-exist", // non-existing
			"-h",
		},
		exps: []Expectation{
			{'x', "TEST"},
			{'h', ""},
		},
	}
	tcase.Test(t);
}

// Dash as argument
func Test_getopt_dash(t *testing.T) {
	Getopt_reset();
	tcase := Test_case{
		cfg_optstr: "x:hX:",
		cfg_longopt:  []Option{
			{"value",      true,  'x'},
			{"method",     true,  'X'},
			{"help",       false, 'h'},
		},
		argv: []string{
			"a.out",   "-x", "-",   "-x-",   "--method",  "-",
		},
		exps: []Expectation{
			{'x', "-"}, // normal `-x -`
			{'x', "-"}, // GNU style `-x-`
			{'X', "-"}, // long `--method -`
		},
	}
	tcase.Test(t);
}

// End of options test
func Test_eof_options(t *testing.T) {
	Getopt_reset();
	const optstr = "x:h"
	lopts := []Option{
		{"method",        true,  'x'},
		{"help",          true,  'h'},
	};

	argv := []string{"a.out", "-h",  "--",  "-x", "val"};
	var idx int;

	idx = Getopt_long (argv, optstr, lopts);
	if idx != 'h' {
		t.Fatalf("test failed")
	}
	idx = Getopt_long (argv, optstr, lopts);
	if idx != -1 {
		t.Fatalf("not correctly ended the getopt")
	}
	if Optind != 2 { // index of '--'
		t.Fatalf("Optind is not set properly")
	}
}
