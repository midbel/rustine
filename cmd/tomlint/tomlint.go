package main

import (
	"flag"
	"log"
	"os"

	"github.com/midbel/rustine/cli/toml"
)

func main() {
	flag.Parse()
	for _, a := range flag.Args() {
		f, err := os.Open(a)
		if err != nil {
			continue
		}
		if err := toml.Valid(f); err != nil {
			log.Printf("%s: %s", f.Name(), err)
		}
		f.Close()
	}
}
