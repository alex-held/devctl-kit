package main

import (
	"fmt"
	"os"

	"github.com/alex-held/devctl-kit/pkg/log"
)

func usage() {
	println(`
		USAGE:
			run-integration-tests [LEVEL] [FORMAT STRING] [OPTIONAL FORMAT ARGS]
		`)
}
func main() {

	if len(os.Args) <= 3 {
		usage()
		os.Exit(1)
	}

	level, err := log.ParseLevel(os.Args[1])
	if err != nil {
		fmt.Printf("tried to parse level '%s' but failed. err=%v\n", os.Args[1], err)
		os.Exit(1)
	}

	formatString := os.Args[2]
	formatArgs := []interface{}{}

	if len(os.Args) > 3 {
		for _, arg := range os.Args[3:] {
			formatArgs = append(formatArgs, arg)
		}
	}

	log.Logf(level, formatString, formatArgs...)
}
