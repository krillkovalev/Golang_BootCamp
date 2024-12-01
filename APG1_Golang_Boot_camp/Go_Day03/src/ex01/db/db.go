package db

// TARGET: вывести страницу со списком имен, адресов и телефонов, чтобы пользователь мог увидеть ее в браузере.

//  TASKS:
// 1)  Вы должны абстрагировать свою базу данных за интерфейсом. Чтобы просто вернуть список записей и иметь возможность листать их/
//   в пакете db все, что связано с базой данных
//       этого интерфейса будет достаточно:
//       type Store interface {
//       // returns a list of items, a total number of hits and (or) an error in case of one
//       GetPlaces(limit int, offset int) ([]types.Place, int, error)
//       }

// 2)  HTTP-приложение должно работать на порту 8888, получая в ответ список ресторанов и обеспечивая простую пагинацию по нему
//         при запросе "http://127.0.0.1:8888/?page=2" (с учетом GET-параметра 'page') вы должны получить такую страницу:
// 3)   Ссылка "Предыдущая" должна исчезнуть на первой странице, а ссылка "Следующая" - на последней.
// 4)  если параметр 'page' указан с неправильным значением (вне [0..last_page] или не является числовым), ваша страница должна вернуть ошибку HTTP 400 и простой текст с описанием ошибки

// curl -s -XGET "http://localhost:9200/places/_search"

import (
	"Interface/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type Store interface {
	// returns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

func GetPlaces(limit int, offset int) ([]types.Place, int, error) {
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

	total := response["count"].(float64)

	esReq := map[string]interface{}{
		"size": limit,
		"from": offset,
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
			if id, ok := sourceMap["id"].(int); ok {
				data.ID = id
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
			if location, ok := sourceMap["location"].(types.GeoPoint); ok {
				data.Location = location
			}
		}
		places = append(places, data)
		fmt.Println(hit)
	}

	return places, int(total), err
}
