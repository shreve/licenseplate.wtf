package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"licenseplate.wtf/db"
	"licenseplate.wtf/model"
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

		db.StartDaemon()
		s := server.NewServer()
		s.ListenAndServe()
	},
}

var name = &cobra.Command{
	Use:   "name",
	Short: "Generate a site name from the input",
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")
		fmt.Println(model.NameHash(input))
	},
}

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Database backup",
}

func init() {
	backup.AddCommand(&cobra.Command{
		Use:     "run",
		Aliases: []string{"", "now"},
		Short:   "Save a backup now",
		Run: func(cmd *cobra.Command, args []string) {
			db.Backup()
		},
	})
	backup.AddCommand(&cobra.Command{
		Use:   "restore",
		Short: "Restore from a backup",
		Run: func(cmd *cobra.Command, args []string) {
			db.Restore()
		},
	})
	rootCmd.AddCommand(serve)
	rootCmd.AddCommand(name)
	rootCmd.AddCommand(backup)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
