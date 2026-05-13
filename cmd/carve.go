package cmd

import (
    "fmt"
    "os"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/carver"
    "github.com/spf13/cobra"
)

var carveCmd = &cobra.Command{
    Use:   "carve [file]",
    Short: "Extract data appended after a PNG IEND chunk",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        title   := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("34"))
        dim     := lipgloss.NewStyle().Foreground(lipgloss.Color("22"))
        green   := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
        yellow  := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
        red     := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
        box     := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("34"))

        result, err := carver.Carve(args[0], "./output/carved")
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        header := title.Render("🧅 unpeeler") + dim.Render("  file carver\n")
        body := fmt.Sprintf("  %-16s %s\n", "File:", args[0])
        body += fmt.Sprintf("  %-16s %s\n", "Source format:", green.Render(string(result.SourceFormat)))

        if !result.HasAppended {
            body += "\n  " + yellow.Render("No appended data found after IEND.")
        } else {
            body += fmt.Sprintf("  %-16s %s\n", "IEND offset:", dim.Render(carver.FormatOffset(result.IENDOffset)))
            body += fmt.Sprintf("  %-16s %s\n", "Appended format:", green.Render(string(result.CarvedFormat)))
            body += fmt.Sprintf("  %-16s %s\n", "Appended size:", formatBytes(result.CarvedSize))
            body += "\n  " + green.Render("✓ Carved to: "+result.CarvedPath)
            body += "\n  " + dim.Render("  Run 'peel' on the carved file to decompress further.")
        }

        if result.SourceFormat != "png" {
            body += "\n\n  " + red.Render("⚠  Extension mismatch or non-PNG source — magic bytes used for detection.")
        }

        fmt.Println(box.Render(header + "\n" + body))
    },
}

func init() {
    rootCmd.AddCommand(carveCmd)
}
