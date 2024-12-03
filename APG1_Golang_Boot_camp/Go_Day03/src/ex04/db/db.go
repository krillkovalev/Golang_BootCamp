package db

// curl -s -XGET "http://localhost:9200/places/_search"

import (
	"Interface/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"github.com/elastic/go-elasticsearch/v8"
)

type Store interface {
	// returns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

func GetPlaces(limit int, lat float64, lon float64) ([]types.Place, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	result, err := es.Count(
		es.Count.WithIndex("places"),
	)
	if err != nil {
		log.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}
	
	esReq := map[string]interface{}{
		"size": limit,
		"sort": []interface{}{
			map[string]interface{}{
				"_geo_distance": map[string]interface{}{
					"location": map[string]float64{
						"lat": lon, //Видимо в индексе перепутаны lat и lon 
						"lon": lat,	
					},
				"order": "asc",
				"unit": "km",
				"mode": "min",
				"distance_type": "arc",
				"ignore_unmapped": true,
				},
			},
		},	
	}
	

	index := "places"

	es.Search()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esReq); err != nil {
		log.Fatal(err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatal(err)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}

	var places []types.Place

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		data := types.Place{}
		if sourceMap, ok := source.(map[string]interface{}); ok {
			if id, ok := sourceMap["id"].(float64); ok {
				data.ID = int(id)
			}
			if name, ok := sourceMap["name"].(string); ok {
				data.Name = name
			}
			if phone, ok := sourceMap["phone"].(string); ok {
				data.Phone = phone
			}
			if name, ok := sourceMap["address"].(string); ok {
				data.Address = name
			}
			if location, ok := sourceMap["location"].(map[string]interface{}); ok {
				lat := location["lat"].(float64)
				data.Location.Lat = lat
				lon := location["lon"].(float64)
				data.Location.Lon = lon
			}
			
		}
		fmt.Println(hit)
		places = append(places, data)
	}

	return places, err
}



