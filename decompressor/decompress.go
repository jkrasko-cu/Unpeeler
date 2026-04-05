package decompressor

import (
    "archive/tar"
    "archive/zip"
    "compress/bzip2"
    "compress/gzip"
    "encoding/base64"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/klauspost/compress/zstd"
    "github.com/ulikunitz/xz"
    "github.com/jkrasko-cu/File-Systems-CLI-Tool/detector"
)

func Decompress(inputPath, outputDir string, format detector.Format) (string, error) {
    os.MkdirAll(outputDir, 0755)
    outPath := filepath.Join(outputDir, "layer")

    f, err := os.Open(inputPath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    switch format {
    case detector.Gzip:
		gr, err := gzip.NewReader(f)
		if err != nil {
			return "", err
		}
		defer gr.Close()
		name := gr.Header.Name
		if name == "" {
			name = "layer.out"
		}
		outPath = filepath.Join(outputDir, name)
		return writeToFile(outPath, gr)
	
	case detector.Bzip2:
		br := bzip2.NewReader(f)
		outPath += ".bz2.out"
		return writeToFile(outPath, br)
	
	case detector.Xz:
		xr, err := xz.NewReader(f)
		if err != nil {
			return "", err
		}
		outPath += ".xz.out"
		return writeToFile(outPath, xr)
	
	case detector.Zstd:
		zr, err := zstd.NewReader(f)
		if err != nil {
			return "", err
		}
		defer zr.Close()
		outPath += ".zst.out"
		return writeToFile(outPath, zr)

    case detector.Zip:
        return extractZip(inputPath, outputDir)

    case detector.Tar:
        return extractTar(f, outputDir)

    case detector.Base64:
        data, err := io.ReadAll(f)
        if err != nil {
            return "", err
        }
        decoded, err := base64.StdEncoding.DecodeString(string(data))
        if err != nil {
            return "", err
        }
        outPath += ".out"
        return outPath, os.WriteFile(outPath, decoded, 0644)
    }

    return "", fmt.Errorf("unsupported format: %s", format)
}

func writeToFile(path string, r io.Reader) (string, error) {
    out, err := os.Create(path)
    if err != nil {
        return "", err
    }
    defer out.Close()
    _, err = io.Copy(out, r)
    return path, err
}

func extractTar(r io.Reader, outputDir string) (string, error) {
    tr := tar.NewReader(r)
    var lastFile string
    for {
        hdr, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return "", err
        }
        target := filepath.Join(outputDir, hdr.Name)
        if hdr.Typeflag == tar.TypeDir {
            os.MkdirAll(target, 0755)
        } else {
            os.MkdirAll(filepath.Dir(target), 0755)
            writeToFile(target, tr)
            lastFile = target
        }
    }
    return lastFile, nil
}

func extractZip(zipPath, outputDir string) (string, error) {
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return "", err
    }
    defer r.Close()
    var lastFile string
    for _, f := range r.File {
        target := filepath.Join(outputDir, f.Name)
        rc, err := f.Open()
        if err != nil {
            return "", err
        }
        writeToFile(target, rc)
        rc.Close()
        lastFile = target
    }
    return lastFile, nil
}