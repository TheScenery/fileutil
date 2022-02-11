package fileutil

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Zip(dst, src string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer zw.Close()

	return filepath.Walk(src, func(fPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		zipFileName, err := filepath.Rel(src, fPath)
		if zipFileName == "" {
			return err
		}
		fh.Name = zipFileName

		if info.IsDir() {
			fh.Name += "/"
		}

		fw, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		fr, err := os.Open(fPath)
		if err != nil {
			return err
		}
		defer fr.Close()

		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}

		return nil

	})
}
