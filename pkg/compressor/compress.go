package compressor

import (
    "archive/tar"
    "archive/zip"
    "compress/gzip"
    "encoding/base64"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/jkrasko-cu/File-Systems-CLI-Tool/pkg/detector"
    "github.com/klauspost/compress/zstd"
    "github.com/ulikunitz/xz"
)

func Compress(inputPath, outputDir string, format detector.Format, depth int) (string, error) {
    os.MkdirAll(outputDir, 0755)

    ext := extFor(format)
    outPath := filepath.Join(outputDir, fmt.Sprintf("layer_%d%s", depth, ext))

    in, err := os.Open(inputPath)
    if err != nil {
        return "", err
    }
    defer in.Close()

    out, err := os.Create(outPath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    switch format {
    case detector.Gzip:
        w := gzip.NewWriter(out)
        _, err = io.Copy(w, in)
        w.Close()
    case detector.Bzip2:
        return "", fmt.Errorf("bzip2 write not supported by stdlib — use gzip, xz, zstd, or zip")
    case detector.Xz:
        w, e := xz.NewWriter(out)
        if e != nil {
            return "", e
        }
        _, err = io.Copy(w, in)
        w.Close()
    case detector.Zstd:
        w, e := zstd.NewWriter(out)
        if e != nil {
            return "", e
        }
        _, err = io.Copy(w, in)
        w.Close()
    case detector.Zip:
        zw := zip.NewWriter(out)
        fw, e := zw.Create(filepath.Base(inputPath))
        if e != nil {
            return "", e
        }
        _, err = io.Copy(fw, in)
        zw.Close()
    case detector.Tar:
        tw := tar.NewWriter(out)
        info, e := os.Stat(inputPath)
        if e != nil {
            return "", e
        }
        hdr, e := tar.FileInfoHeader(info, "")
        if e != nil {
            return "", e
        }
        tw.WriteHeader(hdr)
        _, err = io.Copy(tw, in)
        tw.Close()
    case detector.Base64:
        data, e := io.ReadAll(in)
        if e != nil {
            return "", e
        }
        encoded := base64.StdEncoding.EncodeToString(data)
        _, err = out.WriteString(encoded)
    default:
        return "", fmt.Errorf("unsupported format: %s", format)
    }

    return outPath, err
}

func extFor(f detector.Format) string {
    switch f {
    case detector.Gzip:
        return ".gz"
    case detector.Xz:
        return ".xz"
    case detector.Zstd:
        return ".zst"
    case detector.Zip:
        return ".zip"
    case detector.Tar:
        return ".tar"
    case detector.Base64:
        return ".b64"
    default:
        return ".out"
    }
}
