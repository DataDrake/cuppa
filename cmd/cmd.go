package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

const (
	// GOOD means everthing went perfect
	GOOD = iota
	// USAGE means this command was called incorrectly, and a Usage should be printed
	USAGE = 1
	// FAIL means something went terribly wrong
	FAIL = 2
)

// CMD is a common interface for all CLI commands
type CMD interface {
	Execute() int
	Short() string
	Usage()
}

// Short decription of this utility
var short string

// All of the commands for this application
var subcommands map[string]CMD

// SetShort assigns a short description to this application
func SetShort(s string) {
	short = s
}

// RegisterCMD add a sub-command to this program
func RegisterCMD(name string, c CMD) {
	if subcommands == nil {
		subcommands = make(map[string]CMD)
	}
	subcommands[name] = c
}

// Run finds the appropriate CMD and executes it, or prints the global Usage
func Run() {
	if len(os.Args) < 2 {
		Usage()
	}

	if c := subcommands[os.Args[1]]; c != nil {
		switch c.Execute() {
		case USAGE:
			c.Usage()
		case FAIL:
			os.Exit(1)
		default:
			return
		}
	} else {
		Usage()
	}
}

// Usage prints the usage for this program
func Usage() {
	print("USAGE: " + os.Args[0] + " CMD [OPTIONS] <ARGS>\n\n")
	if len(short) > 0 {
		print("DESCRIPTION: " + short + "\n\n")
	}
	print("COMMANDS:\n\n")
	var keys []string
	i := -1
	for k := range subcommands {
		keys = append(keys, k)
		if len(k) > i {
			i = len(k)
		}
	}
	sort.Strings(keys)
	i += 4
	for _, k := range keys {
		fmt.Printf("%"+strconv.Itoa(i)+"s - %s\n", k, subcommands[k].Short())
	}
	print("\n")
	os.Exit(1)
}
