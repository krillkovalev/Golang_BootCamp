package db

import (
	"Interface/types"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Store interface {
	// returns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

func GetPlaces(limit int, offset int) ([]types.Place, int, error) {
	esReq := map[string]int{
		"size": 10,
	}

	if limit > 0 {
		esReq["limit"] = limit
	}

	if offset > 0 {
		esReq["offset"] = offset
	}

	url := "http://localhost:9200/places/_search"

}

func ReadCsv(filename string) ([]types.Place, error) {
	// Try to open the example.csv file in read-write mode.
	file, err := os.Open(filename)
	// If an error occurs during os.OpenFIle, panic and halt execution.
	if err != nil {
		panic(err)
	}
	// Ensure the file is closed once the function returns
	defer file.Close()

	// Читаем содержимое CSV файла
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Инициализируем slice для хранения записей
	var places []types.Place

	// Проверяем, что файл имеет заголовок
	if len(records) < 1 {
		return nil, fmt.Errorf("empty CSV file or missing header")
	}

	// Обрабатываем каждую строку (начиная с 1, чтобы пропустить заголовок)
	for _, row := range records[1:] {
		// Конвертируем данные из строки в структуру
		id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Fatal(err)
		}
		name := row[1]
		address := row[2]
		phone := row[3]

		// Парсим Latitude и Longitude
		lat, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %v", err)
		}
		lon, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %v", err)
		}

		// Создаем объект Place и заполняем поля
		place := types.Place{
			ID:      id,
			Name:    name,
			Address: address,
			Phone:   phone,
			Location: types.GeoPoint{
				Lat: lat,
				Lon: lon,
			},
		}

		// Добавляем объект Place в slice
		places = append(places, place)
	}

	return places, nil
}

func CreateIndex(indexName string, mapping string, client *elasticsearch.Client) {

	if res, err = client.Indices.Delete([]string{indexName}, client.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res.Body.Close()
	res, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		log.Fatalf("Cannot create index: %s", res)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()
}

func IndexDocument(filename string, indexName string, client *elasticsearch.Client) {

	places, err := ReadCsv(filename)
	if err != nil {
		log.Fatal(err)
	}
	numWorkers := 5
	flushBytes := 5e+6
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,
		Client:        client,
		NumWorkers:    numWorkers,
		FlushBytes:    int(flushBytes),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		log.Print(err)
	}

	start := time.Now().UTC()

	for _, p := range places {
		data, err := json.Marshal(p)
		if err != nil {
			log.Fatal("Cannot encode place: %s", err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(p.ID),

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()

	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}

	log.Println(res)
}
