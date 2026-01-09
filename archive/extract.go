package archive

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/gabriel-vasile/mimetype"
	"github.com/mholt/archiver/v3"
)

// func extract7ZFile(f sevenzip.File) error {
// 	rc, err := f.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer rc.Close()

// 	// Extract the file

// 	return nil
// }

func extract7ZArchive(src string, dest string, overwrite bool) (err error) {
	r, err := sevenzip.OpenReader(src)
	if err != nil {
		return err
	}

	defer func() {
		err = r.Close()
	}()

	for _, f := range r.File {
		filePath := filepath.Join(dest, f.Name)
		// slog.Debug("unzipping file %v", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			slog.Warn("invalid file path")
			return nil
		}
		if f.FileInfo().IsDir() {
			slog.Debug("creating directory...")
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			err = dstFile.Close()
		}()

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			err = fileInArchive.Close()
		}()

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

	}

	return err
}

func Extract(src string, dest string, overwrite bool) error {
	_, err := os.Stat(dest)

	if !os.IsNotExist(err) {
		slog.Warn("Path exists <prompt for overwrite>, skipping")
	} else {
		m, err := mimetype.DetectFile(src)
		if err != nil {
			return err
		}

		if m.Extension() == ".7z" {
			return extract7ZArchive(src, dest, overwrite)
		} else {
			err := archiver.Unarchive(src, dest)
			return err

		}
	}

	return nil
}
