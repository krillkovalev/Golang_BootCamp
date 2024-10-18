package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type DBReader interface {
	Read(filename string, cakeData *Recipes) *Recipes
}

type Ingredient struct {
	Name  string  `json:"ingredient_name" xml:"itemname"`
	Count float64 `json:"ingredient_count,string" xml:"itemcount"`
	Unit  string  `json:"ingredient_unit,omitempty" xml:"itemunit,omitempty"`
}

type Cake struct {
	Name        string       `json:"name" xml:"name"`
	Time        string       `json:"time,omitempty" xml:"stovetime,omitempty"`
	Ingredients []Ingredient `json:"ingredients" xml:"ingredients>item"`
}

type Recipes struct {
	Cakes []Cake `json:"cake" xml:"cake"`
}

type jsonReader struct {
	Filename string
}

type xmlReader struct {
	Filename string
}

func main() {

	filename := flag.String("f", "", "Filename is stores here")

	flag.Parse()

	if strings.Contains(*filename, ".json") {
		jsonFile := &jsonReader{Filename: *filename}
		first_recipe, err := jsonFile.Read()
		if err != nil {
			log.Fatal(err)
		}
		jsonData, err := xml.MarshalIndent(first_recipe, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))

	}

	if strings.Contains(*filename, ".xml") {
		xmlFile := &xmlReader{Filename: *filename}
		second_recipe, err := xmlFile.Read()
		if err != nil {
			log.Fatal(err)
		}
		xmlData, err := json.MarshalIndent(second_recipe, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(xmlData))
	}
}

func (reader *jsonReader) Read() (Recipes, error) {

	recipe := Recipes{}
	file, err := os.Open(reader.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Раскодировать JSON данные в структуру
	if err := json.Unmarshal(data, &recipe); err != nil {
		log.Fatal(err)
	}

	return recipe, err
}

// XMLRead читает XML файл и заполняет структуру
func (reader *xmlReader) Read() (Recipes, error) {

	recipe := Recipes{}
	file, err := os.Open(reader.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Раскодировать XML данные в структуру
	if err := xml.Unmarshal(data, &recipe); err != nil {
		log.Fatal(err)
	}

	return recipe, err
}
