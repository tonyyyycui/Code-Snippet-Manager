package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "snip"}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new snippet",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Adding a snippet...")
		},
	}
	rootCmd.AddCommand(addCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all snippets",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing all snippets...")
		},
	}
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
