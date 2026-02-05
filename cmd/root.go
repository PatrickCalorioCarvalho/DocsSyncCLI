package cmd

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
)

var projectPath string

var rootCmd = &cobra.Command{
	Use:   "docssync",
	Short: "DocsSync CLI - sincronização de documentação técnica",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&projectPath,
		"path",
		"p",
		".",
		"Pasta raiz do projeto",
	)
}
