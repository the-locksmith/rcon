package main

type Config struct {
	ServerAddress  string `json:"server_address"`
	ServerPort     int    `json:"server_port"`
	ServerPassword string `json:"server_password"`

	BotName        string `json:"bot_name"`

	Banlist        []Ban
}

type Ban struct {
	Regex   string
	Period  int
	Message string
}

var config Config
