package flag

import (
	"flag"
	"log"
)

// Require checks that all specified command-line flags are set and not empty.
// It accepts a variable number of flag names as arguments. If any of the provided
// flags are not defined or have empty values, the function logs a fatal error and
// terminates the program. This ensures that all required flags are provided before
// the application continues execution.
func Require(flags ...string) {
	var missingFlags []string

	for _, flagName := range flags {
		flag := flag.Lookup(flagName)
		if flag == nil {
			log.Fatalf("%q is not a valid flag", flagName)
			continue
		}

		if flag.Value.String() == "" {
			missingFlags = append(missingFlags, flagName)
		}
	}

	if len(missingFlags) > 0 {
		log.Fatalf("Missing required flags: %v", missingFlags)
	}
}
