package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	files := os.Args[1:]
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
		out, err := os.Create(file + fmt.Sprintf("_%s", strconv.FormatInt(info.ModTime().Unix(), 10)))
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		err = createArchive(file+fmt.Sprintf("_%s", strconv.FormatInt(info.ModTime().Unix(), 10)), out)
		if err != nil {
			log.Fatal(err)
		}

	}

}

func createArchive(file string, buf io.Writer) error {

	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	err := addToArchive(tw, file)
	if err != nil {
		return err
	}

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = filename

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
