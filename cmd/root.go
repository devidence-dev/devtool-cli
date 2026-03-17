package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "devtool",
	Short: "Developer tooling CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(vscodeCmd)
	rootCmd.AddCommand(vagrantCmd)
}
