rcon
====

[![Build Status](https://travis-ci.org/dieselburner/rcon.svg)](https://travis-ci.org/dieselburner/rcon)
[![License](https://img.shields.io/github/license/dieselburner/rcon.svg)](https://github.com/dieselburner/rcon/blob/master/LICENSE.md)

Source game server RCON command line utility on steroids

Usage
-----

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

Configuration file location
---------------------------
Use default configuration file `~/.rconrc`:

```
rcon status
```

Another option is to specify configuration file via command line or use `RCON_CONF` environment variable, which might be handy in case of maintaining multiple game servers:

```
rcon -config /test/.rconrc status
RCON_CONF=/test/.rconrc rcon status
```

Autoban
-------

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
