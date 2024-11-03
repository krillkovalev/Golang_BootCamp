package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

func main() {

	l := flag.Bool("l", false, "Counting lines in text file")
	m := flag.Bool("m", false, "Counting characters in text file")
	w := flag.Bool("w", false, "Counting words in text file")
	flag.Parse()

	filenames := flag.Args()
	if len(filenames) == 0 {
		fmt.Println("Write filename")
		return
	}

	if flag.NFlag() > 1 {
		log.Fatalf("%s", "Only one flag can be specified at a time")
	}
	var wg sync.WaitGroup
	result := make([]string, len(filenames))

	for i, file := range filenames {
		wg.Add(1)
		switch {
		case *l && !*m && !*w:
			go FileProccess(file, "l", &wg, i, result)
		case *m && !*l && !*w:
			go FileProccess(file, "m", &wg, i, result)

		case *w && !*l && !*m:
			go FileProccess(file, "w", &wg, i, result)
		}

	}

	wg.Wait()

	for _, v := range result {
		fmt.Println(v)
	}

}

func FileProccess(file string, flag string, wg *sync.WaitGroup, i int, result []string) {
	defer wg.Done()

	f, err := os.Open(file)
	if err != nil {
		result[i] = fmt.Sprintf("%s", file)
	}
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	switch flag {
	case "l":
		count, err := LineCounter(f)
		if err != nil {
			log.Fatal(err)
		}
		result[i] = fmt.Sprintf("%d %s", count, file)
	case "m":
		count := CountCharacter(f)
		result[i] = fmt.Sprintf("%d %s", count, file)
	case "w":
		count := CountWords(f)
		result[i] = fmt.Sprintf("%d %s", count, file)
	}
}

func LineCounter(r io.Reader) (int, error) {

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count + 1, nil
}

func CountCharacter(r io.Reader) int {

	count := 0

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		count += utf8.RuneCountInString(line)
	}
	return count

}

func CountWords(r io.Reader) int {
	count := 0

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		count += len(strings.Fields(line))
	}
	return count
}
