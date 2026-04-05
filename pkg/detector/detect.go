package detector

import (
    "os"
    "os/exec"
    "strings"
    "unicode"
)

type Format string

const (
    Gzip    Format = "gzip"
    Bzip2   Format = "bzip2"
    Xz      Format = "xz"
    Zip     Format = "zip"
    Tar     Format = "tar"
    Zstd    Format = "zstd"
    Base64  Format = "base64"
    Unknown Format = "unknown"
)

func Detect(path string) Format {
    out, err := exec.Command("file", "-b", path).Output()
    if err != nil {
        return Unknown
    }
    desc := strings.ToLower(strings.TrimSpace(string(out)))

    switch {
    case strings.Contains(desc, "gzip"):
        return Gzip
    case strings.Contains(desc, "bzip2"):
        return Bzip2
    case strings.Contains(desc, "xz compressed"):
        return Xz
    case strings.Contains(desc, "zip archive"):
        return Zip
    case strings.Contains(desc, "tar archive"):
        return Tar
    case strings.Contains(desc, "zstandard"):
        return Zstd
    case strings.Contains(desc, "base64"):
        return Base64
    default:
        return Unknown
    }
}

func IsHumanReadable(path string) bool {
    out, err := exec.Command("file", "-b", path).Output()
    if err != nil {
        return false
    }
    desc := strings.ToLower(strings.TrimSpace(string(out)))
    if strings.Contains(desc, "ascii") || strings.Contains(desc, "utf-8") || strings.Contains(desc, "text") {
        return true
    }
    data, err := os.ReadFile(path)
    if err != nil || len(data) == 0 {
        return false
    }
    sample := data
    if len(sample) > 512 {
        sample = sample[:512]
    }
    for _, b := range sample {
        if b > 127 || (b < 32 && !unicode.IsSpace(rune(b))) {
            return false
        }
    }
    return true
}
