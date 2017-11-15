package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var DefaultConfig *Config

func init() {
	prgname := strings.ToUpper(filepath.Base(os.Args[0]))
	dirname := os.Getenv(prgname + "_DIR")
	filename := os.Getenv(prgname + "_FILENAME")
	config := os.Getenv(prgname + "_CONFIG")

	var ls []string
	switch runtime.GOOS {
	case "linux":
		ls = []string{"/etc", "/usr/local/etc"}
		if dirname != "" {
			ls = append(ls, dirname)
		}
		if filename == "" {
			filename = prgname
		}
	}
	DefaultConfig = &Config{
		Default:   config,
		Name:      prgname,
		Files:     []string{filename},
		Locations: ls,
	}
}

type Config struct {
	Default   string
	Name      string
	Files     []string
	Locations []string
}

func Configure(v interface{}) error {
	return DefaultConfig.Configure(v)
}

func (c *Config) Configure(v interface{}) error {
	return nil
}
