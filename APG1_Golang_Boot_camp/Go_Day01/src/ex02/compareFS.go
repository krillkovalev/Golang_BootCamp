package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type Snapshot struct {
	filename string
}

func main() {

	old_filename := flag.String("old", "", "Filename of old snapshot")
	new_filename := flag.String("new", "", "Filename of new snapshot")

	flag.Parse()

	old_snapshot := Snapshot{
		filename: *old_filename,
	}
	new_snapshot := Snapshot{
		filename: *new_filename,
	}

	if old_filename == nil || new_filename == nil {
		os.Exit(1)
	}

	ReadCompare(&old_snapshot, &new_snapshot)
}

func ReadCompare(old_snapshot *Snapshot, new_snapshot *Snapshot) error {

	db := make(map[string]bool)
	file, err := os.Open(old_snapshot.filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		db[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file.Close()

	file, err = os.Open(new_snapshot.filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		if !db[scanner.Text()] {
			fmt.Printf("ADDED \"%s\"\n", scanner.Text())
		} else {
			delete(db, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file.Close()

	if len(db) > 0 {
		for k := range db {
			fmt.Printf("REMOVED \"%s\"\n", k)
		}
	}
	return err
}
