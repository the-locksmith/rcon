# rcon

[![Build Status](https://travis-ci.org/dieselburner/rcon.svg)](https://travis-ci.org/dieselburner/rcon)
[![License](https://img.shields.io/github/license/dieselburner/rcon.svg)](https://github.com/dieselburner/rcon/blob/master/LICENSE.md)

Source game server RCON command line utility on steroids

<!-- TOC -->
- [Overview](#overview)
- [Usage](#usage)
- [Configuration](#configuration)
  * [Configuration file format](#configuration-file-format)
  * [Configuration file location](#configuration-file-location)
- [Autoban](#autoban)

## Overview

`rcon` is a RCON command line interface tool for Source engine based games. Apart from obvious use for RCON commands it is capable of automatically banning users based on their user names.

## Usage

```
Usage: rcon [-config file] [-autoban | -autoban-test | -version | command]
  -autoban
        Auto-ban users by their names
  -autoban-test
        Test auto-ban, do not ban anyone
  -config string
        Config file
  -version
        Show version information
```

## Downloads

Latest release is available [here](https://github.com/dieselburner/rcon/releases/latest).

## Configuration

### Configuration file format

Configuration file is a JSON formatted file that contains server information and regex data for autoban feature.

Below is the simplest possible configuration if autoban feature is not used:

```
{
	"server_address"  : "server.example.com",
	"server_port"     : 27015,
	"server_password" : "password"
}
```

Autoban requires some extra configuration:

```
{
	"server_address"  : "server.example.com",
	"server_port"     : 27015,
	"server_password" : "password",

	"bot_name"        : "ban-bot",

	"banlist":
	[
		{
			"regex"   : ".*?banme",
			"period"  : 4320,
			"message" : "*banme* = ban 72h"
		}
	]
}
```

Here, `bot_name` is used as a bot name in generated RCON commands, and is visible to users, while `banlist` is an array of regex rules.

Some technical information for autoban configuration:

- Regex rules are processed in the same order as they appear in the file.
- `period` of 0 used for whitelisting. Because of this no permanent bans are possible via autoban.

Default configuration file is [present in the source code](https://github.com/dieselburner/rcon/blob/master/.rconrc), and contains some predefined set of autoban regex rules.

### Configuration file location

Use default configuration file `~/.rconrc`:

```
rcon status
```

Another option is to specify configuration file via command line or use `RCON_CONF` environment variable, which might be handy in case of maintaining multiple game servers:

```
rcon -config /test/.rconrc status
RCON_CONF=/test/.rconrc rcon status
```

## Autoban

Autoban feature is the primary reason why I developed this tool after I got fed up with advertisement in usernames on my server. Consider using this tool as a cron job if you hate these user names on your server as well:

```
platex | csgo-money.net
ololo twitch.com/ololo
kickwhat.com * myohmy.org
```

Use following to ban/kick these users from your server:

```
rcon -autoban
```

It is possible to test autoban feature without banning anyone, so that no innocent users are banned after configuration adjustment:

```
rcon -autoban-test
```

This will only print matched users along with RCON commands that needs to be performed to ban/kick them, but actual ban/kick will not happen.
