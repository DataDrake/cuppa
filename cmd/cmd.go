package cmd

import (
    "fmt"
)

// CMD is a common interface for all CLI commands
type CMD interface {
    Execute()
    Name() string
    Short() string
    Usage()
}

// All of the commands for this application
var All = []CMD {
    Latest{},
    Quick{},
    Releases{},
}

// Print the usage for this program
func Usage() {
    print("USAGE: cuppa CMD [OPTIONS]\n\n")
    print("COMMANDS:\n\n")
    for _,c := range All {
        fmt.Printf("%12s - %s\n", c.Name(), c.Short())
    }
    print("\n")
}
