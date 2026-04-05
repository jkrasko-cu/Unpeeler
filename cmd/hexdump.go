package cmd

import (
    "fmt"
    "os"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/hexdump"
    "github.com/spf13/cobra"
)

var hexdumpCmd = &cobra.Command{
    Use:   "hexdump [file]",
    Short: "Pretty hex dump of a file with ASCII sidebar",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        result, err := hexdump.Analyze(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        title  := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
        dim    := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        offset := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
        hex    := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
        ascii  := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
        box    := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))

        header := title.Render("🧅 unpeeler") + dim.Render("  hex dump\n")
        body := fmt.Sprintf("  %s\n\n", dim.Render(args[0]))

        for _, line := range result.Lines {
            body += fmt.Sprintf("  %s  %s  %s\n",
                offset.Render(line.Offset),
                hex.Render(line.Hex),
                ascii.Render("|"+line.ASCII+"|"),
            )
        }

        fmt.Println(box.Render(header + "\n" + body))
    },
}

func init() {
    rootCmd.AddCommand(hexdumpCmd)
}
