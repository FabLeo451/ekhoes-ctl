package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var version = "1.0.0"

var port int
var host string
var flagVersion = false

type Config struct {
	URL     string `yaml:"url"`
	Verbose bool   `yaml:"verbose"`
}

type Callback func([]string) int

type Command struct {
	f    Callback
	args string
	help string
}

func confDirExists() (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	path := filepath.Join(home, ".ekhoes")

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

func createEkhoesConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(home, ".ekhoes")

	// 0700 → solo l'utente può accedere
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	configPath := filepath.Join(home, ".ekhoes/conf.yml")

	if _, err := os.Stat(configPath); err == nil {
		return errors.New("il file ~/.ekhoes esiste già")
	} else if !os.IsNotExist(err) {
		return err
	}

	cfg := Config{
		URL:     "https://websocket.ekhoes.com",
		Verbose: false,
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	// 0600 → solo l'utente può leggere/scrivere
	return os.WriteFile(configPath, data, 0600)
}

func main() {
	mapCommands := make(map[string]Command)
	/*
	   mapCommands["programs"] = Command{ ListPrograms, "", "List programs" }
	   mapCommands["compile"] = Command{ CopmileProgram, "ID", "Compile program with the given id" }
	   mapCommands["plugins"] = Command{ ListPlugins, "", "List installed plugins" }
	   mapCommands["install"] = Command{ InstallPlugin, "JARFILE", "Install a plugin from a jar file" }
	*/

	flag.StringVar(&host, "h", "localhost", "Set server host")
	flag.IntVar(&port, "p", 8443, "Set server port")
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

	args := flag.Args()

	if len(args) == 0 {
		os.Exit(0)
	}

	var exitValue int = 0

	if c, found := mapCommands[args[0]]; found {
		exitValue = c.f(args)
	} else {
		fmt.Fprintln(os.Stderr, "Unknown command: ", args[0])
		exitValue = 1
	}

	os.Exit(exitValue)
}
