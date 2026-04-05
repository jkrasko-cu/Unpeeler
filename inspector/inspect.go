package inspector

import (
    "fmt"
    "os"

    "github.com/jkrasko-cu/File-Systems-CLI-Tool/detector"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/entropy"
)

type Result struct {
    Path    string
    Format  detector.Format
    Size    int64
    Magic   string
    Entropy entropy.Result
}

func Inspect(path string) (Result, error) {
    info, err := os.Stat(path)
    if err != nil {
        return Result{}, err
    }

    magic, err := readMagic(path)
    if err != nil {
        return Result{}, err
    }

    format := detector.Detect(path)

    ent, err := entropy.Analyze(path)
    if err != nil {
        return Result{}, err
    }

    return Result{
        Path:    path,
        Format:  format,
        Size:    info.Size(),
        Magic:   magic,
        Entropy: ent,
    }, nil
}

func readMagic(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()

    buf := make([]byte, 8)
    n, err := f.Read(buf)
    if err != nil {
        return "", err
    }

    result := ""
    for i := 0; i < n; i++ {
        if i > 0 {
            result += " "
        }
        result += fmt.Sprintf("%02x", buf[i])
    }
    return result, nil
}
