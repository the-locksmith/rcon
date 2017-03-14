package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/madcitygg/rcon"
)

const version = "1.0.0"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Printf("Usage: %s [-config file] [-autoban | -autoban-test | -version | command]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// Declare variables
	arg_config := new(string)

	// Parse command line arguments
	arg_autoban      := flag.Bool("autoban",      false, "Auto-ban users by their names")
	arg_autoban_test := flag.Bool("autoban-test", false, "Test auto-ban, do not ban anyone")
	arg_version      := flag.Bool("version",      false, "Show version information")
	arg_config        = flag.String("config", "", "Config file")
	flag.Parse()

	// Set config filename if it was not provided
	if *arg_config == "" {
		// Try to extract config filename from RCON_CONF environment variable
		env, set := os.LookupEnv("RCON_CONF")
		if set {
			*arg_config = env
		} else {
			// No config argument & no RCON_CONF are set, use default one
			usr, err := user.Current()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get current user: %v\n", err)
				os.Exit(1)
			}
			*arg_config = fmt.Sprintf("%s/.rconrc", usr.HomeDir)
		}
	}

	// Check if command line arguments are valid
	args := 0
	if *arg_autoban {
		args += 1
	}
	if *arg_autoban_test {
		args += 1
	}
	if *arg_version {
		args += 1
	}
	if len(flag.Args()) > 0 {
		args += 1
	}
	if args != 1 {
		fmt.Fprintln(os.Stderr, "Bad arguments")
		usage()
		os.Exit(1)
	}

	// Show version info before parsing the configuration file
	if *arg_version {
		fmt.Printf("%s %s\n", os.Args[0], version)
		os.Exit(0)
	}

	// Build command from command line argument array
	cmd := strings.Join(flag.Args(), " ")

	// Read configuration file
	data, err := ioutil.ReadFile(*arg_config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read configuration file: %v\n", err)
		os.Exit(1)
	}

	// Parse json data into Config structure
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing configuration file: %v\n", err)
		os.Exit(1)
	}

	// Connect to server
	server, err := rcon.Dial(fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to server: %v\n", err)
		os.Exit(1)
	}
	defer server.Close()

	// Authenticate
	err = server.Authenticate(config.ServerPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not authenticate: %v\n", err)
		os.Exit(1)
	}

	// Perform action
	if *arg_autoban {
		autoban(server, false)
	} else
	if *arg_autoban_test {
		autoban(server, true)
	} else {
		response, err := server.Execute(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server response error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(response.Body)
	}
}
