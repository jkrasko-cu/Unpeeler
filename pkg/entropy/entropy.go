package entropy

import (
    "math"
    "os"
)

type Result struct {
    Score       float64
    Label       string
    ByteCount   int
    FreqTable   [256]int
}

func Analyze(path string) (Result, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Result{}, err
    }

    var freq [256]int
    for _, b := range data {
        freq[b]++
    }

    total := float64(len(data))
    var score float64
    for _, count := range freq {
        if count == 0 {
            continue
        }
        p := float64(count) / total
        score -= p * math.Log2(p)
    }

    return Result{
        Score:     score,
        Label:     label(score),
        ByteCount: len(data),
        FreqTable: freq,
    }, nil
}

func label(score float64) string {
    switch {
    case score < 2.0:
        return "highly structured / repetitive"
    case score < 4.0:
        return "plaintext or source code"
    case score < 6.0:
        return "mixed or structured binary"
    case score < 7.5:
        return "likely compressed"
    default:
        return "likely encrypted or well-compressed"
    }
}
