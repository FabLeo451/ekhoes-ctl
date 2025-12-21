package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

var version = "1.0.0"

var flagVersion = false

type Callback func([]string) error

type Command struct {
	f    Callback
	args string
	help string
}

func main() {
	mapCommands := make(map[string]Command)
	mapCommands["login"] = Command{login, "", "Login"}

	flag.BoolVar(&flagVersion, "v", false, "Show version")

	flag.Usage = func() {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), version)
		fmt.Printf("Usage: %s [options] command [arguments]\n", path.Base(os.Args[0]))

		fmt.Println("\nCommands:")

		for key, element := range mapCommands {
			fmt.Printf("  %-8s %-10s %s\n", key, element.args, element.help)
		}

		fmt.Println("\nOptions:")

		flag.VisitAll(func(f *flag.Flag) {
			a, d := "", ""

			if f.Value.String() != "false" {
				d = "(default: " + f.Value.String() + ")"
				a = "<value>"
			}
			fmt.Printf("  -%s %-10s %s %s\n", f.Name, a, f.Usage, d) // f.Name, f.Value
		})
	}

	flag.Parse()

	if flagVersion {
		fmt.Println(version)
	}

	// Check config
	exists, err := confDirExists()
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		fmt.Println("Initializing...")

		if err = createEkhoesConfig(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Done. Now login to continue")
		os.Exit(0)
	}

	err = loadEkhoesConfig()
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()

	if len(args) == 0 {
		os.Exit(0)
	}

	var exitValue int = 0

	if c, found := mapCommands[args[0]]; found {

		// Checking login
		exists, err = tokenExists()
		if err != nil {
			log.Fatal(err)
		}

		if !exists && args[0] != "login" {
			log.Fatal("Please, login first")
		}

		err = c.f(args)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			exitValue = 1
		}
	} else {
		fmt.Fprintln(os.Stderr, "Unknown command: ", args[0])
		exitValue = 1
	}

	os.Exit(exitValue)
}
