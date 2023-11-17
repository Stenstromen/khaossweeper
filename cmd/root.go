package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stenstromen/khaossweeper/minesweeper"
)

var rootCmd = &cobra.Command{
	Use:   "khaossweeper",
	Long:  "MineSweeper clone written in Go using Fyne",
	Short: "MineSweeper clone written in Go using Fyne",
	RunE:  minesweeper.Minesweeper,
}

func Execute() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: Unable to find the home directory.")
		os.Exit(1)
	}
	kubeconfigDefault := filepath.Join(home, ".kube", "config")

	rootCmd.Flags().BoolP("safe-mode", "s", false, "Show what pod would have been deleted, but don't actually delete it")
	rootCmd.Flags().StringP("kubeconfig", "k", kubeconfigDefault, "Kubeconfig file")
	rootCmd.Flags().StringP("namespace", "n", "default", "A name to say hello to.")
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
