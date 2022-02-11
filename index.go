package fileutil

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path"
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

func readZipFile(dst string, file *zip.File) error {
	fr, err := file.Open()
	if err != nil {
		return err
	}
	defer fr.Close()

	fw, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)
	return err
}

func UnZip(dst, src string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, os.ModeDir)
	if err != nil {
		return err
	}

	for _, file := range zr.File {
		fPath := path.Join(dst, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fPath, os.ModeDir); err != nil {
				return err
			}
			continue
		}

		err := readZipFile(fPath, file)
		if err != nil {
			return err
		}
	}

	return nil
}
