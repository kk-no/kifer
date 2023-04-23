package files

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func UnzipInDirectory(dir string, filter FilterFunc) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".zip" {
			continue
		}
		if err := Unzip(dir, file.Name(), filter); err != nil {
			return err
		}
	}
	return nil
}

func Unzip(dir, zipFile string, filter FilterFunc) error {
	zipFilePath := filepath.Join(dir, zipFile)
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err := CreateFileFromZip(dir, f, filter); err != nil {
			return err
		}
	}
	if err := os.Remove(zipFilePath); err != nil {
		return err
	}
	return nil
}

func CreateFileFromZip(dir string, f *zip.File, filter FilterFunc) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	buf := bytes.NewBuffer(nil)
	r := io.TeeReader(rc, buf)
	if !filter(r) {
		return nil
	}

	to, err := os.OpenFile(filepath.Join(dir, f.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer to.Close()

	from := r
	if buf.Len() != 0 {
		from = buf
	}
	if _, err = io.Copy(to, from); err != nil {
		return err
	}
	return nil
}
