package hexdump

import (
    "fmt"
    "os"
    "strings"
    "unicode"
)

type Result struct {
    Lines []Line
}

type Line struct {
    Offset  string
    Hex     string
    ASCII   string
}

func Analyze(path string) (Result, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Result{}, err
    }

    var lines []Line
    for i := 0; i < len(data); i += 16 {
        end := i + 16
        if end > len(data) {
            end = len(data)
        }
        chunk := data[i:end]

        offset := fmt.Sprintf("%08x", i)

        var hexParts []string
        for _, b := range chunk {
            hexParts = append(hexParts, fmt.Sprintf("%02x", b))
        }
        for len(hexParts) < 16 {
            hexParts = append(hexParts, "  ")
        }
        hex := strings.Join(hexParts[:8], " ") + "  " + strings.Join(hexParts[8:], " ")

        ascii := ""
        for _, b := range chunk {
            if b >= 32 && b < 127 && unicode.IsPrint(rune(b)) {
                ascii += string(b)
            } else {
                ascii += "."
            }
        }

        lines = append(lines, Line{Offset: offset, Hex: hex, ASCII: ascii})
    }

    return Result{Lines: lines}, nil
}
