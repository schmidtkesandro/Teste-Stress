// loadtester/loadtester.go
package loadtester

import (
	"net/http"
	"sync"
)

// ExecuteLoadTest executa o teste de carga e retorna um relatório
type Report struct {
	TotalRequests          int         // Total de requisições enviadas
	SuccessfulRequests     int         // Número de requisições com status
	StatusCodeDistribution map[int]int // Distribuição dos códigos de status HTTP
}

// ExecuteLoadTest performs the load test and returns a report
func ExecuteLoadTest(url string, totalRequests, concurrency int) Report {
	// Inicialização de canais e espera de grupo para controle de concorrência
	var wg sync.WaitGroup
	reqChan := make(chan int, totalRequests)
	resultChan := make(chan int)

	// Envia todas as requisições necessárias para o canal
	for i := 0; i < totalRequests; i++ {
		reqChan <- i
	}
	close(reqChan)

	// Mapa para manter a distribuição dos códigos de status HTTP
	statusCodeDistribution := make(map[int]int)

	// Inicia goroutines para processar as requisições concorrentemente

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Loop para processar todas as requisições do canal
			for range reqChan {
				// Envia requisição HTTP para a URL especificada
				resp, err := http.Get(url)
				//req = req * 1
				if err != nil {
					// Em caso de erro, registra um código de status HTTP interno
					resultChan <- http.StatusInternalServerError
					continue
				}
				// Registra o código de status HTTP recebido
				resultChan <- resp.StatusCode
			}
		}()
	}
	// Fecha o canal de resultados quando todas as goroutines terminarem
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	successfulRequests := 0
	// Processa os resultados recebidos das goroutines
	for code := range resultChan {
		if code == http.StatusOK {
			// Incrementa o contador de requisições bem-sucedidas
			successfulRequests++
		}
		// Atualiza a distribuição dos códigos de status HTTP
		statusCodeDistribution[code]++
	}
	// Retorna um relatório com os resultados do teste de carga
	return Report{
		TotalRequests:          totalRequests,
		SuccessfulRequests:     successfulRequests,
		StatusCodeDistribution: statusCodeDistribution,
	}
}
