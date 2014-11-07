package main

type server struct {
	Address  string   `toml:"address"`
	Channels []string `toml:"channels"`
}

type command struct {
	Name   string `toml:"name"`
	Script string `toml:"script"`
}

type config struct {
	CommandPrefix string    `toml:"command_prefix"`
	Servers       []server  `toml:"server"`
	Commands      []command `toml:"command"`
}
