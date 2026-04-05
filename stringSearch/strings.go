package stringSearch

import (
    "os"
    "strings"
)

const minLength = 4

type Result struct {
    Strings []string
    Count   int
}

func Analyze(path string) (Result, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Result{}, err
    }

    var results []string
    var curr strings.Builder

    for _, b := range data {
        if (b >= 32 && b <= 126) || b == '\t' {
            curr.WriteByte(b)
        } else {
            if curr.Len() >= minLength {
                results = append(results, curr.String())
            }
            curr.Reset()
        }
    }

    if curr.Len() >= minLength {
        results = append(results, curr.String())
    }

    return Result{Strings: results, Count: len(results)}, nil
}
