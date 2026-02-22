// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"

	log "github.com/siamak-amo/v2utils/log"
	pkg "github.com/siamak-amo/v2utils/pkg"
	getopt "github.com/siamak-amo/v2utils/getopt"
)

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
	cmd int                 // CMD_xxx
	urls []string
	configs []string		// file or dir for testing
	output_dir string		// output dir
	in_file string			// input URLs file path
	rm bool					// remove files if broken or invalid
	reverse bool            // print broken configs, not functionals
	verbose bool

	// Internal
	cfg string // config or template file path
	url string

	v2 pkg.V2utils
};

func print_usage() {
	fmt.Fprintf (os.Stderr, "v2utils v%s - xray-core compatible utility\n", Version);
	println (`Usage:  v2utils COMMAND [OPTIONS]

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
        --no-color        disable log color

Test command options:
    -r, --reverse         only print broken configs on stdout
    -R, --rm              to remove broken config files
    -T, --timeout         timeout 2s, 20000ms (default 10s)
    -n, --test-count      number of distinct tests before give up

Examples:
    # run xray by URL:
    $ v2utils run --url 'vless://id@1.2.3.4:1234'

    # test json files and remove broken ones
    $ v2utils test --config /path/to/configs/ --rm

    # convert URLs to json:
    $ cat url.txt | v2utils convert -o /path/to/configs

    # convert outbound of json files to URL:
    $ v2utils convert --config /path/to/configs_dir
`);
}

func (opt *Opt) GetArgs() {
	const optstr = "i:u:f:T:t:o:c:n:RrVvhC"
	lopts := []getopt.Option{
		{"url",           true,  'u'},
		{"config",        true,  'c'},
		{"template",      true,  't'},
		{"output",        true,  'o'},
		{"input",         true,  'i'},

		{"reverse",       false, 'r'},
		{"rm",            false, 'R'},
		{"Timeout",       true,  'T'},
		{"test-count",    true,  'n'},
		{"tc",            true,  'n'},

		{"help",          false, 'h'},
		{"no-color",      false, 'C'},
		{"verbose",       false, 'v'},
		{"version",       false, 'V'},
	}
	argv := os.Args
	for idx := 0; -1 != idx; {
		idx = getopt.Getopt_long (argv, optstr, lopts);

		switch (idx) {
		case 'u':
			opt.urls = append (opt.urls, getopt.Optarg); break;
		case 'c':
			for i := getopt.Optind-1;  i < len(argv) &&
				('-' != argv[i][0] || "-" == argv[i]);  i += 1 {
				opt.configs = append (opt.configs, argv[i]);
			}
			break;
		case 't':
			opt.cfg = getopt.Optarg; break;
		case 'T':
			var e error
			if pkg.TestTimeout, e = time.ParseDuration(getopt.Optarg); nil != e {
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
				pkg.TestCount = count
			}
			break;
		case 'C':
			log.ColorEnabled = false; break;
		case 'V':
			printVersion();
			os.Exit(0);
		case 'h':
			print_usage();
			os.Exit(0);
		}
	}
}

func (opt *Opt) Set2_convert() int {
	if 0 < len(opt.configs) {
		opt.cmd = CMD_CONVERT_CFG;
	} else {
		// Default: converting URLs (stdin / --url)
		opt.cmd = CMD_CONVERT_URL;
	}
	return 0;
}
func (opt *Opt) Set2_test() int {
	if 0 < len(opt.configs) {
		// For Testing config (.json) files
		opt.cmd = CMD_TEST_CFG;
	} else {
		// Default: testing URLs
		opt.cmd = CMD_TEST_URL;
	}
	return 0;
}
func (opt *Opt) Set2_run() int {
	if 0 < len(opt.configs) {
		opt.cmd = CMD_RUN_CFG;
	} else if 0 < len(opt.urls) {
		opt.cmd = CMD_RUN_URL;
	} else {
		log.Warnf("neither URL nor config file is provided, assuming to read URL.\n");
		opt.cmd = CMD_RUN_URL;
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
		case "v", "ver", "version":
			printVersion();
			os.Exit(0);
		default:
			if '-' == argv[1][0] {
				println ("No command is provided.");
			} else {
				fmt.Fprintf (os.Stderr, "Invalid command `%s'.\n", argv[1]);
			}
			println ("Try 'v2utils --help' for more information.");
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

	switch (opt.cmd) {
	case CMD_RUN_URL:
		opt.init_read_url()
		break;

	case CMD_TEST_URL, CMD_CONVERT_URL:
		if "" != opt.output_dir {
			if err := os.MkdirAll(opt.output_dir, 0o755); nil != err {
				log.Errorf ("Could not create dir - %v\n", err);
				return -1
			}
		}
		if 0 < len(opt.urls) {
			opt.init_read_url();
		} else if "" != opt.in_file {
			if e := opt.init_read_file(); nil != e {
				log.Errorf ("%v\n", e);
				return -1
			}
		} else {
			opt.init_read_stdin();
		}
		break;

	case CMD_CONVERT_CFG, CMD_TEST_CFG, CMD_RUN_CFG:
		if 0 == len(opt.configs) {
			opt.init_read_cfg_stdin();
		} else {
			opt.init_read_cfg();
		}
		break;
	}
	return 0;
}

// @return:  0 on success, positive of failure
//    negative on Fatal failures and Exit request
func (opt Opt) Do() int {
	switch (opt.cmd) {
	case CMD_CONVERT_URL:
		if !opt.v2.HasTemplate() {
			if "" != opt.cfg {
				if e := opt.Apply_template(); nil != e {
					log.Errorf("broken or invalid template - %v\n", e);
					return -1;
				}
			} else {
				// No template is provided, using the default one
				opt.Apply_Default_Template();
			}
		}
		if e := opt.Apply_URL(); nil != e {
			log.Warnf("Could not apply URL '%s' - %v\n", opt.url, e);
			return 1;
		}
		if e := opt.MK_josn_output(opt.url); nil != e {
			log.Errorf("IO error: %v\n", e);
			log.Errorf("Fatal error, exiting.\n");
			return -1;
		}
		break;

	case CMD_RUN_URL:
		if !opt.v2.HasTemplate() {
			if e := opt.Init_CFG(); nil != e {
				log.Errorf("Invalid template - %v\n", e)
				return -1;
			}
		}
		if e := opt.Init_Outbound_byURL(); nil != e {
			log.Errorf("Invalid or unsupported URL - %v\n", e)
			return 1;
		}
		if "" == opt.cfg {
			log.Warnf("No template is provided, using the default template: %s\n",
				opt.Get_Default_Template());
		}
		if e := opt.v2.Exec_Xray(); nil != e {
			log.Errorf("Exec xray-core failed - %v\n", e)
			return 1;
		}
		return -1; // The run command, only uses the first provided URL

	case CMD_TEST_URL:
		opt.v2.UnsetTemplate()
		res := opt.Test_URL()
		// Generating json files if applicable
		if res && "" != opt.output_dir {
			if e := opt.MK_josn_output(opt.url); nil != e {
				log.Errorf("IO error: %v\n", e);
				log.Errorf("Fatal error, exiting.\n");
				return -1;
			}
		}
		break;

		// For xxx_CFG commands, @ln is path to a file or `-` for stdin
	case CMD_TEST_CFG:
		res := false
		opt.v2.UnsetTemplate()
		if e := opt.Init_CFG(); nil != e {
			log.Errorf("Loading config file '%s' failed - %v\n", opt.cfg, e)
		} else {
			res = opt.Test_CFG()
		}
		if !res && opt.rm { // We are not in reverse mode here
			if e := os.Remove(opt.cfg); nil != e {
				log.Errorf("Could not remove %s - %v\n", opt.cfg, e)
			} else {
				log.Logf("Broken file %s was removed.\n", opt.cfg)
			}
		}
		break;

	case CMD_RUN_CFG:
		opt.v2.UnsetTemplate()
		if "-" != opt.cfg {
			log.Infof("Loading config: %s\n", opt.cfg)
		}
		if e := opt.Init_CFG(); nil != e {
			log.Errorf("Loading config '%s' failed - %v\n", opt.cfg, e)
			return -1;
		}
		if !opt.v2.HasInboundConfig() {
			log.Warnf(
				"No 'inbounds' section found in '%s', using the default template: %s\n",
				opt.cfg, opt.Get_Default_Template(),
			)
			opt.v2.SetDefaultInboundConfig();
		}
		if e := opt.v2.Exec_Xray(); nil != e {
			log.Errorf("Exec xray-core failed - %v\n", e)
			return -1;
		}
		return -1; // RUN_CFG only uses the first input

	case CMD_CONVERT_CFG:
		opt.v2.UnsetTemplate()
		if e := opt.Init_CFG(); nil != e {
			log.Errorf("Loading config '%s' failed - %v\n", opt.cfg, e)
			return 1;
		}
		res, e := opt.v2.Convert_conf2url();
		if nil != e {
			log.Warnf ("Converting '%s' to URL failed - %v\n", opt.cfg, e);
		} else {
			fmt.Println(res);
		}
		break;
	}

	return 0;
}

func init_opt() (opt *Opt) {
	opt = &Opt{};
	opt.GetArgs();
	if ret := opt.HandleArgs(); ret < 0 {
		os.Exit (-ret);
	}
	if ret := opt.Init(); ret < 0 {
		os.Exit (-ret);
	}
	return;
}

// main loop of v2utils program (blocking)
func main_loop(opt *Opt) {
	for ;; {
		if EOF := opt.GetInput(); true == EOF {
			break;
		}
		if opt.Do() < 0 {
			break;
		}
	}
}
