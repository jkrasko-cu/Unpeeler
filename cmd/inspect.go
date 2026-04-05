package cmd

import (
    "fmt"
    "math"
    "os"
    "strings"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/inspector"
    "github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
    Use:   "inspect [file]",
    Short: "Inspect a file: format, size, magic bytes, and entropy",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        result, err := inspector.Inspect(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        title  := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
        dim    := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        label  := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
        box    := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))
        scored := scoreStyle(result.Entropy.Score)

        header := title.Render("🧅 unpeeler") + dim.Render("  file inspection\n")

        bar := renderInspectBar(result.Entropy.Score)

        formatStr := string(result.Format)
        if result.Format == "unknown" {
            formatStr = dim.Render("unknown (may be plaintext)")
        } else {
            formatStr = label.Render(formatStr)
        }

        body := fmt.Sprintf("  %-10s %s\n", "File:", args[0])
        body += fmt.Sprintf("  %-10s %s\n", "Format:", formatStr)
        body += fmt.Sprintf("  %-10s %s\n", "Size:", formatBytes(int(result.Size)))
        body += fmt.Sprintf("  %-10s %s\n", "Magic:", dim.Render(result.Magic))
        body += fmt.Sprintf("  %-10s %s  %s / 8.0\n", "Entropy:", bar, scored.Render(fmt.Sprintf("%.2f", result.Entropy.Score)))
        body += fmt.Sprintf("  %-10s %s\n", "Verdict:", scored.Render(result.Entropy.Label))

        fmt.Println(box.Render(header + "\n" + body))
    },
}

func renderInspectBar(score float64) string {
    const width = 20
    filled := int(math.Round(score / 8.0 * width))
    if filled > width {
        filled = width
    }
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    return scoreStyle(score).Render(bar)
}

func init() {
    rootCmd.AddCommand(inspectCmd)
}
