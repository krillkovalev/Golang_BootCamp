package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func main() {

	a := flag.Bool("a", false, "Multiple files for archiving")
	flag.Parse()

	var wg sync.WaitGroup
	var files []string
	if *a {
		files = os.Args[3:]
	} else {
		files = append(files, os.Args[1])
	}

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			if *a {
				createArchive(f, os.Args[2], true)
			} else {
				createArchive(f, filepath.Dir(os.Args[1]), false)
			}

		}(file)

	}
	wg.Wait()

}

func createArchive(file string, path string, flag bool) error {
	info, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	timestamp := strconv.FormatInt(info.ModTime().Unix(), 10)
	base_filename := filepath.Base(file)
	out, err := os.Create(getOutputPath(path, base_filename, timestamp, flag))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	addToArchive(tw, file)
	return nil
}

func getOutputPath(path string, base_filename string, timestamp string, flag bool) string {
	if flag {
		return path + "/" + strings.Replace(base_filename, ".", "_", -1) + fmt.Sprintf("_%s", timestamp) + ".tar.gz"

	}
	return path + "/" + strings.Replace(base_filename, ".", "_", -1) + fmt.Sprintf("_%s", timestamp) + ".tar.gz"
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
