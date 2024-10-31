package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var path string
	d := flag.Bool("d", false, "Flag for directory searching")
	file := flag.Bool("f", false, "Flag for file searching")
	ext := flag.String("ext", "", "File extension filter (works ONLY when -f is specified)")
	sl := flag.Bool("sl", false, "Flag for symlink searching")
	flag.Parse()

	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}
	if !*d && !*file && !*sl {
		searchPath(path, "", "")
	}
	if *file && len(*ext) > 0 {
		searchPath(path, "ext", *ext)
		*file = false
	}
	for *d || *file || *sl {
		switch {
		case *d:
			searchPath(path, "d", "")
			*d = false
		case *file:
			searchPath(path, "f", "")
			*file = false
		case *sl:
			searchPath(path, "sl", "")
			*sl = false
		}
	}
}

func searchPath(path string, flag string, ext string) {
	seen := make(map[string]bool, 1000)
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			fileinfo, err := os.Stat(path)
			if err != nil && (flag == "sl" || flag == "") {
				fmt.Printf("%s -> [broken]\n", filepath.Dir(path))
				return nil
			}
			switch flag {
			case "d":
				dirpath := filepath.Dir(path)
				if _, exists := seen[dirpath]; exists {
					return nil
				}
				seen[dirpath] = true
				if info.IsDir() {
					fmt.Println(path)
				}
			case "f":
				PrintFlagF(path, fileinfo)
			case "ext":
				FlagWithExt(path, ext, fileinfo)
			case "sl":
				PrintFlagSL(path)
			default:
				if info.IsDir() {
					fmt.Println(path)
				}
				PrintFlagF(path, fileinfo)
				PrintFlagSL(path)
			}

			return nil
		})
	if err != nil {
		os.Exit(1)
	}
}

func PrintFlagF(path string, fileinfo os.FileInfo) bool {
	if fileinfo == nil {
		return false
	}
	if !fileinfo.IsDir() && !isSymlink(path) {
		fmt.Println(path)
	}
	return true
}

func FlagWithExt(path string, ext string, fileinfo os.FileInfo) bool {
	if fileinfo == nil {
		return false
	}
	if !fileinfo.IsDir() && !isSymlink(path) {
		if filepath.Ext(path) == ("." + ext) {
			fmt.Println(path)
		}
	}
	return true
}

func PrintFlagSL(path string) {
	if isSymlink(path) {
		dst, err := os.Readlink(path)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%s -> %s\n", filepath.Dir(path), dst)
		}
	}
}

func isSymlink(path string) bool {
	stat, err := os.Lstat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return stat.Mode()&os.ModeSymlink == os.ModeSymlink // symlink check here

}
