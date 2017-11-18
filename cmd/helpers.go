package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runcmd(name string, args ...string) {

	var outBuff bytes.Buffer
	var errBuff bytes.Buffer
	gcmd := exec.Command(name, args...)
	gcmd.Stdout = &outBuff
	gcmd.Stderr = &errBuff
	//nolint
	gcmd.Run()

	if outBuff.Len() > 0 {
		log.Println(outBuff.String())
	}
	if errBuff.Len() > 0 {
		log.Println(errBuff.String())
	}

}

func globExt(dir string, ext string) ([]string, error) {

	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func glob(dir string, name string) ([]string, error) {

	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Base(path) == name {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
