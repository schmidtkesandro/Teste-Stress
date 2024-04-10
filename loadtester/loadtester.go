// loadtester/loadtester.go
package loadtester

import (
	"net/http"
	"sync"
)

// Report struct to hold load test results
type Report struct {
	TotalRequests          int
	SuccessfulRequests     int
	StatusCodeDistribution map[int]int
}

// ExecuteLoadTest performs the load test and returns a report
func ExecuteLoadTest(url string, totalRequests, concurrency int) Report {
	var wg sync.WaitGroup
	reqChan := make(chan int, totalRequests)
	resultChan := make(chan int)

	for i := 0; i < totalRequests; i++ {
		reqChan <- i
	}
	close(reqChan)

	statusCodeDistribution := make(map[int]int)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for req := range reqChan {
				resp, err := http.Get(url)
				req = req * 1
				if err != nil {
					resultChan <- http.StatusInternalServerError
					continue
				}
				resultChan <- resp.StatusCode
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	successfulRequests := 0

	for code := range resultChan {
		if code == http.StatusOK {
			successfulRequests++
		}
		statusCodeDistribution[code]++
	}

	return Report{
		TotalRequests:          totalRequests,
		SuccessfulRequests:     successfulRequests,
		StatusCodeDistribution: statusCodeDistribution,
	}
}
