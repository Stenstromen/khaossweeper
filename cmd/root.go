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
	Short: "KhaosSweeper: A thrilling Kubernetes minesweeper game.",
	Long: `KhaosSweeper brings a unique twist to the classic minesweeper game. 
Inspired by Chaos Monkey and Kube DOOM, KhaosSweeper integrates with your Kubernetes cluster to randomly kill pods when you hit a mine. 
This game is designed not just for fun, but also as a novel way to test the resilience and fault tolerance of your Kubernetes setup. 
Experience the thrill of navigating through a minefield where each wrong move could bring down a pod! 
KhaosSweeper is built using Go, FyneV2, and k8s.io/client-go.`,
	RunE: minesweeper.Minesweeper,
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
