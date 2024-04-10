// main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sandroschmidtke/teste-stress/loadtester"
)

func main() {
	var url string
	var requests int
	var concurrency int

	flag.StringVar(&url, "url", "", "URL do serviço a ser testado")
	flag.IntVar(&requests, "requests", 0, "Número total de requests")
	flag.IntVar(&concurrency, "concurrency", 0, "Número de chamadas simultâneas")
	flag.Parse()

	if url == "" || requests <= 0 || concurrency <= 0 {
		fmt.Println("Por favor, forneça a URL do serviço, o número total de requests e o nível de concorrência no formato:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Exemplo:")
		fmt.Println("go run main.go --url=http://google.com --requests=1000 --concurrency=10")
		fmt.Println("")
		os.Exit(1)
	}

	start := time.Now()
	report := loadtester.ExecuteLoadTest(url, requests, concurrency)
	elapsed := time.Since(start)

	fmt.Printf("\nTempo total gasto: %v\n", elapsed)
	fmt.Printf("Quantidade total de requests: %d\n", report.TotalRequests)
	fmt.Printf("Quantidade de requests com status 200: %d\n", report.SuccessfulRequests)
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range report.StatusCodeDistribution {
		fmt.Printf("- Status %d: %d\n", status, count)
	}
}
