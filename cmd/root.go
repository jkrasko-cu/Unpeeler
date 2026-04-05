package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/jkrasko-cu/File-Systems-CLI-Tool/ui"
)

var outputPath string

var rootCmd = &cobra.Command{
	Use: "Unpeeler [file]",
	Short: "Decompresses nested archives",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		m := ui.NewModel(inputFile, outputPath)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "./output", "Output directory")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}