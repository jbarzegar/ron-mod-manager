package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/gabriel-vasile/mimetype"
	"github.com/mholt/archiver/v3"
)

func extract7ZFile(f sevenzip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Extract the file

	return nil
}

func extract7ZArchive(src string, dest string, overwrite bool) error {
	r, err := sevenzip.OpenReader(src)

	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(dest, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return nil
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func ExtractArchive(src string, dest string, overwrite bool) error {

	_, err := os.Stat(dest)

	if !os.IsNotExist(err) {
		fmt.Println("Path exists", "<prompt for overwrite>, skipping")
	} else {
		m, err := mimetype.DetectFile(src)

		if err != nil {
			return err
		}

		if m.Extension() == ".7z" {
			extract7ZArchive(src, dest, overwrite)
		} else {
			err := archiver.Unarchive(src, dest)

			return err

		}
	}

	return nil
}
