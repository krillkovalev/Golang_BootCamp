package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/lizrice/secure-connections/utils"
)

type Candy struct {
	Money			int		`json:"money"`
	CandyType 		string	`json:"candyType"`
	CandyCount 		int		`json:"candyCount"`
}


func main() {

	candy_name := flag.String("k", "", "select candyType")
	candy_count := flag.Int("c", 0, "count of candy to buy")
	amount := flag.Int("m", 0, "amount of money you gave to machine")

	flag.Parse()

	client := getClient()

	order := Candy{
		Money: 			*amount,
		CandyType: 		*candy_name,
		CandyCount: 	*candy_count,
		
	}

	marshalled, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("impossible to marshall")
	}

	request, err := http.NewRequest("-XPOST", "https://candy.ltd/buy_candy:3333", bytes.NewReader(marshalled))
	if err != nil {
		return
	}

	// добавляем заголовки
    request.Header.Add("Content-Type", "application/json")  
  
    resp, err := client.Do(request) 
    if err != nil { 
        fmt.Println(err) 
        return
    } 

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}


	fmt.Printf("Status: %s Body: %s\n", resp.Status, string(body))

}

func getClient() *http.Client {
	cp := x509.NewCertPool()
	data, _ := os.ReadFile("../ca/minica.pem")
	cp.AppendCertsFromPEM(data)


	config := &tls.Config{
		RootCAs: 		cp,
		GetClientCertificate: utils.ClientCertReqFunc("cert.pem", "key.pem"),
		VerifyPeerCertificate: utils.CertificateChains,
	} 

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}

	return &client
}

