package toml

import (
	"strings"
	"testing"
)

type user struct {
	Name   string `toml:"name"`
	Passwd string `toml:"passwd"`
}

type pool struct {
	Proxy string `toml:"proxy"`
	User  user   `toml:"auth"`
	Hosts []conn `toml:"databases"`
}

type conn struct {
	Server  string   `toml:"server"`
	Ports   []uint16 `toml:"ports"`
	Limit   uint16   `toml:"limit"`
	Enabled bool     `toml:"enabled"`
	User    user     `toml:"auth"`
}

func TestDecodeNestedTables(t *testing.T) {
	s := `
title = "TOML parser test"

[pool]
proxy = "192.168.1.124"
auth = {name = "midbel", passwd = "midbeltoml101"}

[[pool.databases]]
server = "192.168.1.1"
ports = [8001, 8002, 8003]
limit = 10

[[pool.databases]]
server = "192.168.1.1"
ports = [8001, 8002, 8003]
limit = 10
auth = {name = "midbel", passwd = "tomlrules101"}
  `
	c := struct {
		Title    string `toml:"title"`
		Pool pool `toml:"pool"`
	}{}
	if err := NewDecoder(strings.NewReader(s)).Decode(&c); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeTableArray(t *testing.T) {
	s := `
title = "TOML parser test"

[[database]]
server = "192.168.1.1"
ports = [8001, 8002, 8003]
limit = 10

[[database]]
server = "192.168.1.1"
ports = [8001, 8002, 8003]
limit = 10
auth = {name = "midbel", passwd = "tomlrules101"}
  `
	c := struct {
		Title    string `toml:"title"`
		Settings []conn `toml:"database"`
	}{}
	if err := NewDecoder(strings.NewReader(s)).Decode(&c); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeSimple(t *testing.T) {
	s := `
title = "TOML parser test"

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
limit = 5000
enabled = true
  `
	c := struct {
		Title    string `toml:"title"`
		Settings conn   `toml:"database"`
	}{}
	if err := NewDecoder(strings.NewReader(s)).Decode(&c); err != nil {
		t.Fatal(err)
	}
}