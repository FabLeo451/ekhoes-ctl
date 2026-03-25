package main

import (
    "fmt"
    "os"

    appcmd "ekhoes-ctl/cmd"
    "ekhoes-ctl/config"

    "github.com/spf13/cobra"
)

var version = "1.0.0"

var (
    urlFlag        string
    tokenFlag      string
    showVersionFlg bool
)

func main() {
    rootCmd := &cobra.Command{
        Use:           "ekhoes-ctl",
        Short:         "ekhoes-ctl - CLI for managing Ekhoes resources",
        Long:          "ekhoes-ctl is a command line interface by which users can perform operation on Ekhoes environment, monitor resources and much more.",
        Version:       version,
        SilenceUsage:  true,
        SilenceErrors: false,
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            // Gestione esplicita -v/--version (oltre a --version di Cobra)
            /*if showVersionFlg {
                fmt.Println(version)
                os.Exit(0)
            }*/

            // 1) Controllo/inizializzazione config
            exists, err := config.ConfDirExists()
            if err != nil {
                return err
            }
            if !exists {
                fmt.Println("Initializing...")
                if err := config.CreateEkhoesConfig(); err != nil {
                    return err
                }
                fmt.Println("Done. Now login to continue")
                os.Exit(0)
            }

            if err := config.LoadEkhoesConfig(); err != nil {
                return err
            }

            // 2) Override URL se passato via flag
            if cmd.Flags().Changed("url") && urlFlag != "" {
                config.Conf.URL = urlFlag
            }

            // 3) Override token se passato via flag
            if cmd.Flags().Changed("token") && tokenFlag != "" {
                appcmd.SetToken(tokenFlag)
            }
			
			if cmd.Name() != "completion" {
				return nil
			}

            // 4) Controllo token per tutti i comandi tranne "login"
            //    In un PersistentPreRunE del root, `cmd` è il comando specifico invocato.
            if cmd.Name() != "login" {
                if tokenFlag == "" { // se non arriva da flag, prova da storage
                    tok, err := appcmd.GetToken()
                    if err != nil {
                        return err
                    }
                    if tok == "" {
                        return fmt.Errorf("Please, login first")
                    }
                }
            }
            return nil
        },
    }

    // Flag globali (ereditati da tutti i sottocomandi)
    rootCmd.PersistentFlags().StringVarP(&urlFlag, "url", "u", "", "Server URL")
    rootCmd.PersistentFlags().StringVarP(&tokenFlag, "token", "t", "", "Authentication token")
    rootCmd.PersistentFlags().BoolVarP(&showVersionFlg, "version", "v", false, "Show version")

    // --- Comandi applicativi ---

    rootCmd.AddCommand(&cobra.Command{
        Use:   "login",
        Short: "Login",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.Login(append([]string{"login"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "logout",
        Short: "Delete authentication token",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.Logout(append([]string{"logout"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "sessions",
        Short: "Retrieve sessions",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.GetSessions(append([]string{"sessions"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "kill <session_id>",
        Short: "Kill a session",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.KillSession(append([]string{"kill"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "killall",
        Short: "Kill all sessions",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.KillAllSessions(append([]string{"killall"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "connections",
        Short: "Retrieve websocket connections",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.GetWebsocketConnections(append([]string{"connections"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "system",
        Short: "Retrieve system information",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.GetSystemInfo(append([]string{"system"}, args...))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "top",
        Short: "Retrieve running processes",
        RunE: func(cmd *cobra.Command, args []string) error {
            return appcmd.TopCpuProcesses(append([]string{"top"}, args...))
        },
    })

    // --- Comando completion ---
    completionCmd := &cobra.Command{
        Use:       "completion [bash|zsh|fish|powershell]",
        Short:     "Genera script di completamento per la shell",
        Long:      "Genera lo script di completamento per bash, zsh, fish o powershell.",
        ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
        Args:      cobra.ExactValidArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            sh := args[0]
            switch sh {
            case "bash":
                return rootCmd.GenBashCompletion(os.Stdout)
            case "zsh":
                // Suggerimento: per zsh è utile anche la 'compinit' nell'rc
                return rootCmd.GenZshCompletion(os.Stdout)
            case "fish":
                // true = include descriptions
                return rootCmd.GenFishCompletion(os.Stdout, true)
            case "powershell":
                // Disponibile nelle versioni recenti di cobra
                // In alternativa, rootCmd.GenPowerShellCompletion(os.Stdout)
                return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
            default:
                return fmt.Errorf("shell non supportata: %s", sh)
            }
        },
    }
    rootCmd.AddCommand(completionCmd)

    // Esecuzione
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}