/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CarlosEduardoAD/broadcast-server/internal/server"
	"github.com/CarlosEduardoAD/broadcast-server/internal/utils"
	"github.com/spf13/cobra"
)

var (
	port = 8080
	srv  = server.NewServer(port)
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server [subcommand]",
	Short: "Commands to manage the websocket server",
	Long:  `The server command allows the user to start and stop it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "start":

			cUser, err := os.UserHomeDir()

			if err != nil {
				log.Fatal(err)
			}

			filePath, err := filepath.Abs(cUser + "/projects/broadcast-server/broadcast-server.pid")
			pidFile, err := os.Open(filePath)

			if err != nil {
				ok := os.IsNotExist(err)

				if !ok {
					panic(err)
				}
			}

			if os.Getenv("IS_DAEMON") != "1" {

				if pidFile != nil {
					os.Exit(0)
				}

				pid, err := utils.Fork()
				if err != nil {
					fmt.Printf("Unable to fork process")
					os.Exit(1)
				}

				fmt.Printf("Start child - pid = %d\n", pid)
				log.Println("Starting server...")
				os.Exit(0) // parent must exit after spawning child
			}

			srv.Connect()
		case "stop":
			cUser, err := os.UserHomeDir()

			if err != nil {
				log.Fatal(err)
			}
			log.Println(cUser)

			filePath, err := filepath.Abs(cUser + "/projects/broadcast-server/broadcast-server.pid")
			pidFile, err := os.ReadFile(filePath)

			pdrInt, err := strconv.Atoi(string(pidFile))

			if err != nil {
				panic(err)
			}

			process, err := os.FindProcess(pdrInt)
			if err != nil {
				log.Fatalf("failed to find process: %v", err)
			}

			err = process.Kill()
			if err != nil {
				log.Fatalf("failed to kill process with pid %s: %v", pdrInt, err)
			}
			os.Remove(filePath)

			fmt.Printf("Server with pid %s stopped successfully\n", pdrInt)

		default:
			fmt.Println("Unknown subcommand. Use 'start' or 'stop'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
