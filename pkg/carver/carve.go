package carver

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"

    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/detector"
)

var iendMarker = []byte{0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

type Result struct {
    SourceFormat  detector.Format
    CarvedFormat  detector.Format
    IENDOffset    int
    CarvedSize    int
    CarvedPath    string
    HasAppended   bool
}

func Carve(path, outputDir string) (Result, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Result{}, err
    }

    sourceFormat := detector.Detect(path)

    pos := bytes.Index(data, iendMarker)
    if pos == -1 {
        return Result{SourceFormat: sourceFormat, HasAppended: false}, nil
    }

    afterIEND := pos + len(iendMarker)
    if afterIEND >= len(data) {
        return Result{SourceFormat: sourceFormat, IENDOffset: pos, HasAppended: false}, nil
    }

    carved := data[afterIEND:]

    // write carved bytes to temp file so detector can read it
    os.MkdirAll(outputDir, 0755)
    tmpPath := filepath.Join(outputDir, "carved.tmp")
    if err := os.WriteFile(tmpPath, carved, 0644); err != nil {
        return Result{}, err
    }

    carvedFormat := detector.Detect(tmpPath)

    ext := extFor(carvedFormat)
    finalPath := filepath.Join(outputDir, "carved"+ext)
    if err := os.Rename(tmpPath, finalPath); err != nil {
        return Result{}, err
    }

    return Result{
        SourceFormat: sourceFormat,
        CarvedFormat: carvedFormat,
        IENDOffset:   pos,
        CarvedSize:   len(carved),
        CarvedPath:   finalPath,
        HasAppended:  true,
    }, nil
}

func extFor(f detector.Format) string {
    switch f {
    case detector.Gzip:   return ".gz"
    case detector.Zip:    return ".zip"
    case detector.Bzip2:  return ".bz2"
    case detector.Xz:     return ".xz"
    case detector.Zstd:   return ".zst"
    case detector.Tar:    return ".tar"
    default:              return ".bin"
    }
}

func FormatOffset(offset int) string {
    return fmt.Sprintf("0x%X", offset)
}
