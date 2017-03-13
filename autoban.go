package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/madcitygg/rcon"
)

type user_data struct {
	id   string
	name string
}

// Convert multiline server response string to array of lines
func string_to_array(s string) []string {
	var array []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		array = append(array, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse server response: %v\n", err)
		os.Exit(1)
	}

	return array
}

// Check user status string against ban list
func check_match(banlist []Ban, s string) (*user_data, *Ban) {
	// Setup user status regex and look for the match against input
	re := regexp.MustCompile("(?i).*?\"(.*?)\" (.*?) (.*?) (.*?) (.*?) (.*?) (.*?) (.*?):(.*?)$")
	match := re.FindStringSubmatch(s)

	// If match is successful, it will contain following data
	// 0: full match
	// 1: name
	// 2: steam id
	// 3: connected
	// 4: ping
	// 5: loss
	// 6: state
	// 7: rate
	// 8: ip
	// 9: port

	// Matched, this is the user status string
	if match != nil {
		// Save relevant user data
		user := new(user_data)
		user.id   = match[2]
		user.name = match[1]

		// Iterate through ban records
		for _, ban := range banlist {
			// Setup ban regex and look for the match against username
			re = regexp.MustCompile("(?i)" + ban.Regex)
			match = re.FindStringSubmatch(user.name)
			if match != nil {
				// Username matches one of the banlist records, return both user and ban record
				return user, &ban
			}
		}

		// Username is clean, no ban matches found, return user only
		return user, nil
	}

	// Could not extract user id & name from the string - not a user status string
	return nil, nil
}

// Auto-ban users by their names
func autoban(server *rcon.RCON, test bool) {
	// Get user list, "users" command does not return Steam ID's - use "status" instead
	response, err := server.Execute("status")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server `status` command error: %v\n", err)
		os.Exit(1)
	}

	// Convert response to array of strings
	array := string_to_array(response.Body)

	// Check every response line
	for _, line := range array {
		// Parse the line and check if username matches any of the banlist records
		user, ban := check_match(config.Banlist, line)
		// Match found
		if user != nil && ban != nil {
			if ban.Period > 0 {
				// Ban > 0, ban!
				fmt.Printf("Ban: id=%s name=\"%s\" period=%d, message=\"%s: %s\"\n",
				           user.id, user.name, ban.Period, config.BotName, ban.Message)
				if !test {
					// Execute rcon command: banid <period> <userid>
					cmd := fmt.Sprintf("banid %d %s", ban.Period, user.id)
					fmt.Printf("  > %s\n", cmd)
					response, err := server.Execute(cmd)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Server `banid` command error: %v\n", err)
						os.Exit(1)
					}
					fmt.Printf("  < %s\n", response.Body)
					// Execute rcon command: kickid <userid> <message>
					cmd = fmt.Sprintf("kickid %s %s: %s", user.id, config.BotName, ban.Message)
					fmt.Printf("  > %s\n", cmd)
					response, err = server.Execute(cmd)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Server `kickid` command error: %v\n", err)
						os.Exit(1)
					}
					fmt.Printf("  < %s\n", response.Body)
				}
			}
		}
	}
}
