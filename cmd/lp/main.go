package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"licenseplate.wtf/server"
)

var dev = true

var rootCmd = &cobra.Command{
	Use:   "lp",
	Short: "licenseplate.wtf management tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		if dev {
			// Set up SASS watcher
			go func() {
				cmd := exec.Command(
					"/usr/bin/sass",
					"--style=compressed",
					"--no-source-map",
					"--no-error-css",
					"--watch",
					"server/static/app.sass",
					"server/static/app.css",
				)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()

				cmd = exec.Command(
					"memcached",
					"--port",
					"11211",
				)
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()
		}

		s := server.NewServer()
		s.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serve)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
