package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type DBReader interface {
	Read(cakeData *Recipes) (Recipes, error)
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

	old_filename := flag.String("old", "", "Filename of old database")
	new_filename := flag.String("new", "", "Filename of new database")

	var old_recipes, new_recipes Recipes

	flag.Parse()

	oldFileReader := getFileReader(*old_filename)
	newFileReader := getFileReader(*new_filename)

	if oldFileReader == nil || newFileReader == nil {
		os.Exit(1)
	}

	_, err := oldFileReader.Read(&old_recipes)
	if err != nil {
		log.Fatal(err)
	}

	_, err = newFileReader.Read(&new_recipes)
	if err != nil {
		log.Fatal(err)
	}
	compareCakeNames(old_recipes, new_recipes)
	compareCookingTime(old_recipes, new_recipes)
	compareIngredients(old_recipes, new_recipes)
	Units(old_recipes, new_recipes)

}

func (reader *jsonReader) Read(cakeData *Recipes) (Recipes, error) {

	file, err := os.Open(reader.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, cakeData); err != nil {
		log.Fatal(err)
	}

	return *cakeData, err
}

func (reader *xmlReader) Read(cakeData *Recipes) (Recipes, error) {

	file, err := os.Open(reader.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := xml.Unmarshal(data, cakeData); err != nil {
		log.Fatal(err)
	}

	return *cakeData, err
}

func getFileReader(filename string) DBReader {
	switch filepath.Ext(filename) {
	case ".json":
		return &jsonReader{Filename: filename}
	case ".xml":
		return &xmlReader{Filename: filename}
	default:
		return nil
	}
}

func compareCakeNames(old_recipes Recipes, new_recipes Recipes) {
	oldcakeNames := make([]string, len(old_recipes.Cakes))
	newcakeNames := make([]string, len(new_recipes.Cakes))
	for i, cake := range old_recipes.Cakes {
		oldcakeNames[i] = cake.Name
		for i, cake := range new_recipes.Cakes {
			newcakeNames[i] = cake.Name
		}
	}
	slices.SortFunc(oldcakeNames, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})
	slices.SortFunc(newcakeNames, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})
	for i := 0; i < len(newcakeNames); i++ {
		if newcakeNames[i] == oldcakeNames[i] {
			continue
		} else {
			fmt.Printf("ADDED cake: \"%s\"\n", newcakeNames[i])
			fmt.Printf("REMOVED cake: \"%s\"\n", oldcakeNames[i])
		}
	}

}

func compareCookingTime(old_recipes Recipes, new_recipes Recipes) {

	oldcakeNames := make(map[string]string)
	newcakeNames := make(map[string]string)

	for _, cake := range old_recipes.Cakes {
		oldcakeNames[cake.Name] = cake.Time
		for _, cake := range new_recipes.Cakes {
			newcakeNames[cake.Name] = cake.Time
		}
	}

	for k, v := range oldcakeNames {
		for key, value := range newcakeNames {
			if _, exists := newcakeNames[k]; exists {
				if k == key && v != value {
					fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", key, value, v)
				}
			}
		}
	}

}

func compareIngredients(old_recipes Recipes, new_recipes Recipes) {

	for _, original := range old_recipes.Cakes {
		for _, stolen := range new_recipes.Cakes {
			if original.Name == stolen.Name {
				oldcakeNames := make(map[string]bool)
				newcakeNames := make(map[string]bool)

				for _, ingredient := range original.Ingredients {
					oldcakeNames[ingredient.Name] = true
				}

				for _, ingredient := range stolen.Ingredients {
					newcakeNames[ingredient.Name] = true
				}

				for _, ingredient := range stolen.Ingredients {
					if !oldcakeNames[ingredient.Name] {
						fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", ingredient.Name, stolen.Name)
					}
				}

				for _, ingredient := range original.Ingredients {
					if !newcakeNames[ingredient.Name] {
						fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", ingredient.Name, original.Name)
					}
				}
			}
		}
	}

}

func Units(old_recipes Recipes, new_recipes Recipes) {

	for _, original := range old_recipes.Cakes {
		for _, stolen := range new_recipes.Cakes {
			if original.Name == stolen.Name {
				oldCountNames := make(map[string]float64)
				newCountNames := make(map[string]float64)
				oldUnitNames := make(map[string]string)
				newUnitNames := make(map[string]string)
				oldUnits := make(map[string]bool)
				newUnits := make(map[string]bool)

				for _, ingredient := range original.Ingredients {
					oldCountNames[ingredient.Name] = ingredient.Count
					oldUnitNames[ingredient.Name] = ingredient.Unit
					if _, exists := oldUnits[ingredient.Unit]; !exists {
						oldUnits[ingredient.Unit] = true
					}
				}

				for _, ingredient := range stolen.Ingredients {
					newCountNames[ingredient.Name] = ingredient.Count
					newUnitNames[ingredient.Name] = ingredient.Unit
					if _, exists := newUnits[ingredient.Unit]; !exists {
						newUnits[ingredient.Unit] = true
					}
				}

				for name, unit := range oldUnitNames {
					for n, u := range newUnitNames {
						if name == n && unit != u && unit != "" && u != "" {
							fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\n", n, original.Name, u, unit)
						}
						if name == n && u == "" && len(unit) > 0 {
							defer fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n", unit, n, original.Name)
						}

					}
				}

				for name, count := range oldCountNames {
					for n, cnt := range newCountNames {
						if name == n && count != cnt && oldUnitNames[name] == newUnitNames[n] {
							fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%0.f\" instead of \"%0.f\"\n", name, original.Name, cnt, count)
						}
					}
				}
			}
		}
	}

}
