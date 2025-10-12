package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	UserAgent string `json:"userAgent"`
	Token     string `json:"token"`
	Message   string `json:"message"`
	Receiver  string `json:"receiver"`
}

type LoadTestResult struct {
	TotalRequests  int
	SuccessCount   int
	FailCount      int
	TotalTime      time.Duration
	AverageTime    time.Duration
	RequestsPerSec float64
}

func main() {
	fmt.Println("🔥 Load Testing the Consumer API")
	fmt.Println("This will test your current browser pool implementation via HTTP API")

	// Test different concurrent loads
	testCases := []int{100}

	for _, concurrency := range testCases {
		fmt.Printf("\n📡 Testing with %d concurrent requests...\n", concurrency)
		result := loadTest(concurrency, "http://localhost:8000/api/v1/jobs")
		printResult(concurrency, result)

		// Wait between tests to let the system recover
		time.Sleep(2 * time.Second)
	}
}

func loadTest(concurrency int, endpoint string) LoadTestResult {
	var wg sync.WaitGroup
	results := make(chan bool, concurrency)

	// Test message
	message := Message{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		Token:     "test-token",
		Message:   "Test message from load test",
		Receiver:  "https://linkedin.com/in/test-user",
	}

	start := time.Now()

	// Launch concurrent requests
	for i := range concurrency {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			success := sendRequest(endpoint, message)
			results <- success

			if success {
				fmt.Printf("✅ Request %d completed\n", index+1)
			} else {
				fmt.Printf("❌ Request %d failed\n", index+1)
			}
		}(i)
	}

	wg.Wait()
	close(results)

	totalTime := time.Since(start)

	// Count results
	successCount := 0
	failCount := 0

	for success := range results {
		if success {
			successCount++
		} else {
			failCount++
		}
	}

	avgTime := totalTime / time.Duration(concurrency)
	reqPerSec := float64(concurrency) / totalTime.Seconds()

	return LoadTestResult{
		TotalRequests:  concurrency,
		SuccessCount:   successCount,
		FailCount:      failCount,
		TotalTime:      totalTime,
		AverageTime:    avgTime,
		RequestsPerSec: reqPerSec,
	}
}

func sendRequest(endpoint string, message Message) bool {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return false
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Read response body to ensure complete request
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func printResult(concurrency int, result LoadTestResult) {
	fmt.Printf("📊 Results for %d concurrent requests:\n", concurrency)
	fmt.Printf("   ✅ Success: %d/%d (%.1f%%)\n",
		result.SuccessCount,
		result.TotalRequests,
		float64(result.SuccessCount)/float64(result.TotalRequests)*100)
	fmt.Printf("   ⏱️  Total Time: %v\n", result.TotalTime.Round(time.Millisecond))
	fmt.Printf("   📈 Requests/sec: %.2f\n", result.RequestsPerSec)
	fmt.Printf("   ⚡ Avg Response: %v\n", result.AverageTime.Round(time.Millisecond))
}
