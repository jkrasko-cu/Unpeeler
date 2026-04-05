package cmd

import (
    "fmt"
    "math"
    "os"
    "strings"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/entropy"
    "github.com/spf13/cobra"
)

var entropyCmd = &cobra.Command{
    Use:   "entropy [file]",
    Short: "Calculate Shannon entropy of a file",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        result, err := entropy.Analyze(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        title  := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
        dim    := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        box    := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))
        scored := scoreStyle(result.Score)

        header := title.Render("🧅 unpeeler") + dim.Render("  entropy analysis\n")

        bar := renderBar(result.Score)

        body := fmt.Sprintf("  File:     %s\n", args[0])
        body += fmt.Sprintf("  Size:     %s\n", formatBytes(result.ByteCount))
        body += fmt.Sprintf("  Entropy:  %s  %s / 8.0\n", bar, scored.Render(fmt.Sprintf("%.2f", result.Score)))
        body += fmt.Sprintf("  Verdict:  %s\n", scored.Render(result.Label))

        fmt.Println(box.Render(header + "\n" + body))
    },
}

func renderBar(score float64) string {
    const width = 20
    filled := int(math.Round(score / 8.0 * width))
    if filled > width {
        filled = width
    }
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    return scoreStyle(score).Render(bar)
}

func scoreStyle(score float64) lipgloss.Style {
    switch {
    case score < 4.0:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // green
    case score < 6.5:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // yellow
    default:
        return lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // red
    }
}

func formatBytes(n int) string {
    switch {
    case n < 1024:
        return fmt.Sprintf("%d B", n)
    case n < 1024*1024:
        return fmt.Sprintf("%.1f KB", float64(n)/1024)
    default:
        return fmt.Sprintf("%.1f MB", float64(n)/(1024*1024))
    }
}

func init() {
    rootCmd.AddCommand(entropyCmd)
}
