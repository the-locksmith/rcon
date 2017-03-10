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

type Config struct {
	ServerAddress  string `json:"server_address"`
	ServerPort     int    `json:"server_port"`
	ServerPassword string `json:"server_password"`
}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Printf("Usage: %s [-config file] command\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments
	file := new(string)
	file = flag.String("config", "", "Config file")
	flag.Parse()

	// Set config filename if it was not provided
	if *file == "" {
		// Try to extract config filename from RCON_CONF environment variable
		env, set := os.LookupEnv("RCON_CONF")
		if set {
			*file = env
		} else {
			// No config argument & no RCON_CONF are set, use default one
			usr, err := user.Current()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get current user: %v\n", err)
				os.Exit(1)
			}
			*file = fmt.Sprintf("%s/.rconrc", usr.HomeDir)
		}
	}

	// Check if command line arguments are valid
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Command missing")
		usage()
		os.Exit(1)
	}

	// Build command from command line argument array
	cmd := strings.Join(flag.Args(), " ")

	// Read configuration file
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read configuration file: %v\n", err)
		os.Exit(1)
	}

	// Parse json data into Config structure
	var config Config
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

	// Execute command
	response, err := server.Execute(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server response error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(response.Body)
}
