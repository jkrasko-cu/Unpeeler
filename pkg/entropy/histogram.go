package entropy

import (
    "math"
    "os"
)

type HistogramResult struct {
    Chunks []ChunkResult
}

type ChunkResult struct {
    Offset int
    Score  float64
    Label  string
}

func Histogram(path string, chunkSize int) (HistogramResult, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return HistogramResult{}, err
    }

    var chunks []ChunkResult
    for i := 0; i < len(data); i += chunkSize {
        end := i + chunkSize
        if end > len(data) {
            end = len(data)
        }
        chunk := data[i:end]

        var freq [256]int
        for _, b := range chunk {
            freq[b]++
        }

        total := float64(len(chunk))
        var score float64
        for _, count := range freq {
            if count == 0 {
                continue
            }
            p := float64(count) / total
            score -= p * math.Log2(p)
        }

        chunks = append(chunks, ChunkResult{
            Offset: i,
            Score:  score,
            Label:  label(score),
        })
    }

    return HistogramResult{Chunks: chunks}, nil
}
