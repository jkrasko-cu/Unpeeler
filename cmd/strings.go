package cmd

import (
    "fmt"
    "os"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/stringSearch"
    "github.com/spf13/cobra"
)

var stringsCmd = &cobra.Command{
    Use:   "strings [file]",
    Short: "Extract printable strings from a binary file",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        result, err := stringSearch.Analyze(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        title  := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
        dim    := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        green  := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
        box    := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))

        header := title.Render("🧅 unpeeler") + dim.Render("  strings extraction\n")
        body := fmt.Sprintf("  File:     %s\n", args[0])
        body += fmt.Sprintf("  Found:    %s\n\n", green.Render(fmt.Sprintf("%d strings", result.Count)))
        for _, s := range result.Strings {
            body += fmt.Sprintf("  %s\n", dim.Render(s))
        }

        fmt.Println(box.Render(header + "\n" + body))
    },
}

func init() {
    rootCmd.AddCommand(stringsCmd)
}
