package cmd

import (
    "fmt"
    "math"
    "os"
    "strings"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/entropy"
    "github.com/spf13/cobra"
)

var histogramCmd = &cobra.Command{
    Use:   "histogram [file]",
    Short: "Show per-region entropy across a file",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        result, err := entropy.Histogram(args[0], 256)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
        dim   := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        box   := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))

        header := title.Render("🧅 unpeeler") + dim.Render("  entropy histogram\n")
        body := fmt.Sprintf("  %-10s %-22s %s\n\n", "Offset", "Entropy", "Verdict")

        for _, chunk := range result.Chunks {
            bar := renderHistBar(chunk.Score)
            body += fmt.Sprintf("  %08x  %s %s\n",
                chunk.Offset,
                bar,
                dim.Render(fmt.Sprintf("%.2f  %s", chunk.Score, chunk.Label)),
            )
        }

        fmt.Println(box.Render(header + "\n" + body))
    },
}

func renderHistBar(score float64) string {
    const width = 20
    filled := int(math.Round(score / 8.0 * float64(width)))
    if filled > width {
        filled = width
    }
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    return scoreStyle(score).Render(bar)
}

func init() {
    rootCmd.AddCommand(histogramCmd)
}
