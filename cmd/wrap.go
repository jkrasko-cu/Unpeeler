package cmd

import (
    "fmt"
    "os"
    "strings"

    "github.com/charmbracelet/lipgloss"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/compressor"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/detector"
    "github.com/spf13/cobra"
)

var (
    layers    string
    wrapOut   string
)

var wrapCmd = &cobra.Command{
    Use:   "wrap [file]",
    Short: "Wrap a file in multiple compression layers",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        inputFile := args[0]
        formats := parseFormats(layers)
        if len(formats) == 0 {
            fmt.Println("No valid formats provided. Use --layers gzip,xz,zstd etc.")
            os.Exit(1)
        }

        done := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
        dim  := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
        box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))

        header := title.Render("🧅 unpeeler") + dim.Render("  wrap mode\n")
        body := ""

        current := inputFile
        for i, f := range formats {
            out, err := compressor.Compress(current, wrapOut, f, i+1)
            if err != nil {
                fmt.Println("Error:", err)
                os.Exit(1)
            }
            body += fmt.Sprintf("  %s  Layer %-2d  [%-8s]\n", done, i+1, string(f))
            current = out
        }

        footer := "\n" + dim.Render(fmt.Sprintf("  💾 Output: %s", current))
        fmt.Println(box.Render(header + "\n" + body + footer))
    },
}

func parseFormats(s string) []detector.Format {
    var result []detector.Format
    for _, part := range strings.Split(s, ",") {
        part = strings.TrimSpace(strings.ToLower(part))
        switch part {
        case "gzip", "gz":
            result = append(result, detector.Gzip)
        case "xz":
            result = append(result, detector.Xz)
        case "zstd", "zst":
            result = append(result, detector.Zstd)
        case "zip":
            result = append(result, detector.Zip)
        case "tar":
            result = append(result, detector.Tar)
        case "base64", "b64":
            result = append(result, detector.Base64)
        }
    }
    return result
}

func init() {
    wrapCmd.Flags().StringVarP(&layers, "layers", "l", "gzip", "Comma-separated compression layers (e.g. gzip,xz,zstd)")
    wrapCmd.Flags().StringVarP(&wrapOut, "output", "o", "./wrapped", "Output directory")
    rootCmd.AddCommand(wrapCmd)
}
