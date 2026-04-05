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
    Use:   "unpeeler [file]",
    Short: "🧅 A recursive file forensics toolkit",
    Args:  cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            // Launch interactive TUI
            p := tea.NewProgram(ui.NewTUIModel(), tea.WithAltScreen())
            if _, err := p.Run(); err != nil {
                fmt.Println("Error:", err)
                os.Exit(1)
            }
            return
        }
        // Direct peel mode
        m := ui.NewModel(args[0], outputPath)
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
