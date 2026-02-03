/* V2utils provides xray-core compatible utilities
   Copyright 2025-2026 Ahmad <edu.siamak@gmail.com>

   V2utils is free software: you can redistribute it and/or modify it
   under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   V2utils is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
   See the GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"
	"strconv"

	log "github.com/siamak-amo/v2utils/log"
	utils "github.com/siamak-amo/v2utils/utils"
	getopt "github.com/siamak-amo/v2utils/getopt"
)

const Version = "1.3";

const (
	CMD_CONVERT_URL int = iota
	CMD_CONVERT_CFG
	CMD_TEST_URL
	CMD_TEST_CFG
	CMD_RUN_URL
	CMD_RUN_CFG
) // commands

type Opt struct {
	// User options
	Cmd int                 // CMD_xxx
	url string
	configs string			// file or dir for testing
	output_dir string		// output dir
	template_file string	// template file path
	in_file string			// input URLs file path
	rm bool					// remove files if broken or invalid
	reverse bool            // print broken configs, not functionals
	verbose bool

	scanner *bufio.Scanner
	GetInput func() (string, bool)

	V2 utils.V2utils
};

func (opt *Opt) GetArgs() {
	const optstr = "u:f:T:t:o:c:n:Rrvh"
	lopts := []getopt.Option{
		{"url",           true,  'u'},
		{"config",        true,  'c'},
		{"template",      true,  't'},
		{"Timeout",       true,  'T'},
		{"output",        true,  'o'},
		{"input",         true,  'i'},
		{"reverse",       false, 'r'},
		{"rm",            false, 'R'},
		{"verbose",       false, 'v'},
		{"help",          false, 'h'},
		{"test-count",    true,  'n'},
		{"tc",            true,  'n'},
	}
	argv := os.Args
	for idx := 0; -1 != idx; {
		idx = getopt.Getopt_long (argv, optstr, lopts);

		switch (idx) {
		case 'u':
			opt.url = getopt.Optarg; break;
		case 'c':
			opt.configs = getopt.Optarg; break;
		case 't':
			opt.template_file = getopt.Optarg; break;
		case 'T':
			var e error
			if utils.TestTimeout, e = time.ParseDuration(getopt.Optarg); nil != e {
				log.Errorf("set timeout option failed - %v\n", e);
			}
			break;
		case 'i':
			opt.in_file = getopt.Optarg; break;
		case 'o':
			opt.output_dir = getopt.Optarg; break;
		case 'R':
			opt.rm = true; break;
		case 'r':
			opt.reverse = true; break;
		case 'v':
			opt.verbose = true; break;
		case 'n':
			if count, err := strconv.Atoi(getopt.Optarg); nil == err && count > 0 {
				utils.TestCount = count
			}
			break;
		case 'h':
			fmt.Fprintf(os.Stderr, `v2utils v%s - xray-core compatible utility
Usage:  v2utils COMMAND [OPTIONS]

COMMAND:
      Run:  to execute Xray based on the given configuration
     Test:  to test the current configuration has internet access
  Convert:  to convert the current configuration to a different format

OPTIONS:
    -u, --url             VPN url (e.g. vless:// trojan://)
    -c, --config          path to config file or folder
    -t, --template        path to template file
                          (for Run and Convert commands)
    -i, --input           path to input URL file
    -o, --output          path to output folder
    -v, --verbose         verbose

Test command options:
    -r, --reverse         only print broken configs on stdout
    -R, --rm              to remove broken config files
    -T, --timeout         timeout 20s, 20000ms (default 10s)
    -n, --test-count      number of distinct tests before give up

Examples:
    # run xray by URL:
    $ v2utils run --url vless://id@1.2.3.4:1234

    # test json files and remove broken ones
    $ v2utils test --config /path/to/configs/ --rm

    # convert URLs to json:
    $ cat url.txt | v2utils convert -o /path/to/configs

    # convert outbound of json files to URL:
    $ v2utils convert --config /path/to/configs_dir
`, Version);
			os.Exit(0);
		}
	}
}

func (opt *Opt) Set2_convert() int {
	if "" != opt.configs {
		opt.Cmd = CMD_CONVERT_CFG;
	} else {
		// Default: converting URLs (stdin / --url)
		opt.Cmd = CMD_CONVERT_URL;
	}
	return 0;
}
func (opt *Opt) Set2_test() int {
	if "" != opt.configs {
		// For Testing config (.json) files
		opt.Cmd = CMD_TEST_CFG;
	} else {
		// Default: testing URLs
		opt.Cmd = CMD_TEST_URL;
	}
	return 0;
}
func (opt *Opt) Set2_run() int {
	if "" != opt.configs {
		opt.Cmd = CMD_RUN_CFG;
	} else if "" != opt.url {
		opt.Cmd = CMD_RUN_URL;
	} else {
		log.Warnf("neither URL nor config file is provided, assuming to read URL.\n");
		opt.Cmd = CMD_RUN_URL;
	}
	return 0;
}

// returns negative on fatal failures
func (opt *Opt) HandleArgs() int {
	argv := os.Args
	switch (argv[0][strings.LastIndexByte(argv[0], '/')+1:]) {
	case "v2test":
		return opt.Set2_test();
	case "v2convert", "v2conv":
		return opt.Set2_convert();
	case "v2run", "v2ray", "xray", "xrun":
		return opt.Set2_run();
	default:
		if len(argv) < 2 {
			fmt.Fprintln(os.Stderr, "error:  missing COMMAND")
			fmt.Fprintln(os.Stderr, "usage:  v2utils [COMMAND] [OPTIONS]")
			return -1
		}
		switch (argv[1]) {
		case "convert","CONVERT", "conv", "c","C":
			return opt.Set2_convert();
		case "test","Test","TEST", "t","T":
			return opt.Set2_test();
		case "run","Run","RUN", "r","R":
			return opt.Set2_run();
		default:
			println ("Invalid command.");
			return -1
		}
	}
	return 0;
}

// returns negative on fatal failures
func (opt *Opt) Init() int {
	if opt.verbose {
		log.LogLevel = log.Verbose;
	} else {
		log.LogLevel = log.Warning; // Default level
	}
	if opt.rm && opt.reverse {
		log.Errorf("cannot pass --rm and --reverse options together\n");
		return -1
	}

	switch (opt.Cmd) {
	case CMD_RUN_URL:
		opt.Set_rd_url()
		break;

	case CMD_TEST_URL, CMD_CONVERT_URL:
		if "" != opt.output_dir {
			if err := os.MkdirAll(opt.output_dir, 0o755); nil != err {
				log.Errorf ("Could not create dir - %v\n", err);
				return -1
			}
		}
		if "" != opt.url {
			opt.Set_rd_url();
		} else if "" != opt.in_file {
			if e := opt.Set_rd_file(); nil != e {
				log.Errorf ("%v\n", e);
				return -1
			}
		} else {
			opt.Set_rd_stdin();
		}
		break;

	case CMD_CONVERT_CFG, CMD_TEST_CFG, CMD_RUN_CFG:
		if "-" == opt.configs {
			opt.Set_rd_cfg_stdin();
		} else {
			opt.Set_rd_cfg();
		}
		break;
	}
	return 0;
}

// @return:  0 on success, positive of failure
//    negative on Fatal failures and Exit request
func (opt Opt) Do(ln string) int {
	switch (opt.Cmd) {
	case CMD_CONVERT_URL:
		if !opt.V2.HasTemplate() {
			if "" != opt.template_file {
				if e := opt.V2.Apply_template(opt.template_file); nil != e {
					log.Errorf("broken or invalid template - %v\n", e);
					return -1;
				}
			} else {
				// No template is provided, using the default one
				opt.Apply_Default_Template();
			}
		}
		if e := opt.V2.Apply_URL(ln); nil != e {
			log.Warnf("Could not apply URL '%s' - %v\n", ln, e);
			return 1;
		}
		if e := opt.MK_josn_output(ln); nil != e {
			log.Errorf("IO error: %v\n", e);
			log.Errorf("Fatal error, exiting.\n");
			return -1;
		}
		break;

	case CMD_RUN_URL:
		if !opt.V2.HasTemplate() {
			if e := opt.Init_CFG(); nil != e {
				log.Errorf("Invalid template - %v\n", e)
				return -1;
			}
		}
		if e := opt.V2.Init_Outbound_byURL(ln); nil != e {
			log.Errorf("Invalid or unsupported URL - %v\n", e)
			return 1;
		}
		if "" == opt.template_file {
			log.Warnf("No template is provided, using the default template: %s\n",
				opt.Get_Default_Template());
		}
		if e := opt.V2.Exec_Xray(); nil != e {
			log.Errorf("Exec xray-core failed - %v\n", e)
			return 1;
		}
		return -1; // The run command, only uses the first provided URL

	case CMD_TEST_URL:
		opt.V2.UnsetTemplate()
		result, duration := opt.V2.Test_URL(ln)
		if ! opt.reverse {
			if result {
				fmt.Println(ln)
				if opt.verbose {
					log.Infof("`%s` OK (%d ms).\n", ln, duration)
				}
			} else {
				log.Warnf("Broken URL '%s'\n", ln);
			}
		} else { // Only print broken urls
			if ! result {
				fmt.Println(ln);
			} else if opt.verbose {
				log.Infof("`%s` OK (%d ms).\n", ln, duration)
			}
		}
		// Generating json files if applicable
		if result && "" != opt.output_dir {
			if e := opt.MK_josn_output(ln); nil != e {
				log.Errorf("IO error: %v\n", e);
				log.Errorf("Fatal error, exiting.\n");
				return -1;
			}
		}
		break;

		// For xxx_CFG commands, @ln is path to a file or `-` for stdin
	case CMD_TEST_CFG:
		res := false
		var duration int64
		opt.template_file = ln
		opt.V2.UnsetTemplate()
		if e := opt.Init_CFG(); nil != e {
			log.Errorf("Loading config file '%s' failed - %v\n", ln, e)
		} else {
			res, duration = opt.V2.Test_CFG(ln)
		}
		if ! opt.reverse {
			if res {
				fmt.Println(ln);
				if opt.verbose {
					log.Infof ("config file '%s':  OK (%d ms).\n", ln, duration);
				}
			} else {
				log.Warnf("config file '%s' is broken\n", ln)
			}
		} else { // Only print broken configs
			if !res {
				fmt.Println(ln);
			} else {
				log.Infof("config file '%s':  OK (%d ms).\n", ln, duration)
			}
		}

		if !res && opt.rm { // We are not in reverse mode here
			if e := os.Remove(ln); nil != e {
				log.Errorf("Could not remove %s - %v\n", ln, e)
			} else {
				log.Logf("Broken file %s was removed.\n", ln)
			}
		}
		break;

	case CMD_RUN_CFG:
		opt.template_file = ln
		opt.V2.UnsetTemplate()
		if "-" != ln {
			log.Infof("Loading config: %s\n", ln)
		}
		if e := opt.Init_CFG(); nil != e {
			log.Errorf("Loading config '%s' failed - %v\n", ln, e)
			return -1;
		}
		if e := opt.V2.Exec_Xray(); nil != e {
			log.Errorf("Exec xray-core failed - %v\n", e)
			return -1;
		}
		return -1; // RUN_CFG only uses the first input

	case CMD_CONVERT_CFG:
		opt.template_file = ln
		opt.V2.UnsetTemplate()
		if e := opt.Init_CFG(); nil != e {
			log.Errorf("Loading config '%s' failed - %v\n", ln, e)
			return 1;
		}
		res, e := opt.V2.Convert_conf2url();
		if nil != e {
			log.Warnf ("Converting '%s' to URL failed - %v\n", ln, e);
		} else {
			fmt.Println(res);
		}
		break;
	}

	return 0;
}

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
	var ln string
	for EOF := false; true != EOF; {
		ln, EOF = opt.GetInput();
		if len(ln) == 0 || ln[0] <= ' ' || ln[0] == '#' {
			continue;
		}
		if opt.Do(ln) < 0 {
			break;
		}
	}
}
