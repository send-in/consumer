package main

import (
	queue "consumer/internal/queue"
	"consumer/pkg/browser"
	logger "consumer/pkg/log"
	"consumer/pkg/worker"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type LoadTestResult struct {
	TotalJobs       int
	CompletedJobs   int
	FailedJobs      int
	ProcessingTime  time.Duration
	JobsPerSecond   float64
	AvgTimePerJob   time.Duration
	BrowserPoolSize int
	WorkerThreads   int
}

func main() {
	// Start logger (same as your main.go)
	logger.Start()

	fmt.Println("🔥 Direct Browser Pool & Worker Load Test")
	fmt.Println("Testing YOUR EXACT implementation with 100 jobs...")
	fmt.Println(strings.Repeat("=", 60))

	// Test with your current configuration
	result := testYourExactImplementation(100, 3, 3) // 100 jobs, 3 browsers, 3 workers
	printLoadTestResult(result)
}

func testYourExactImplementation(numJobs, poolSize, maxWorkers int) LoadTestResult {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create browser pool EXACTLY like your main.go
	fmt.Printf("🧠 Creating browser pool with %d browsers (your exact implementation)...\n", poolSize)
	browserPool, err := browser.CreatePool(poolSize, ctx)
	if err != nil {
		logger.Fatal(err, "Failed to start browsers")
		return LoadTestResult{}
	}
	defer browserPool.Close()

	// Create jobs channel EXACTLY like RabbitMQ would
	jobs := make(chan amqp.Delivery, numJobs)

	// Create test messages in YOUR format
	testMessage := queue.Message{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		Token:     "test-token",
		Message:   "Load test message",
		Receiver:  "https://linkedin.com/in/test-user",
	}

	// Convert to JSON (same as RabbitMQ would send)
	messageBytes, err := json.Marshal(testMessage)
	if err != nil {
		logger.Fatal(err, "Failed to marshal test message")
		return LoadTestResult{}
	}

	// Create mock AMQP deliveries with completion tracking
	fmt.Printf("📦 Creating %d test jobs in YOUR message format...\n", numJobs)

	var completedJobs, failedJobs int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create deliveries that track completion
	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		delivery := createMockDelivery(
			fmt.Sprintf("load-test-job-%d", i+1),
			messageBytes,
			&wg,
			&mu,
			&completedJobs,
			&failedJobs,
		)
		jobs <- delivery
	}
	close(jobs)

	// Start timing
	fmt.Printf("⚙️ Starting YOUR worker.Factory with %d threads...\n", maxWorkers)
	startTime := time.Now()

	// Start YOUR EXACT worker.Factory function
	go worker.Factory(jobs, browserPool, maxWorkers, 3, ctx)

	// Wait for all jobs to complete
	wg.Wait()
	totalTime := time.Since(startTime)

	// Calculate metrics
	mu.Lock()
	final_completed := completedJobs
	final_failed := failedJobs
	mu.Unlock()

	jobsPerSecond := float64(final_completed) / totalTime.Seconds()
	avgTimePerJob := time.Duration(0)
	if final_completed > 0 {
		avgTimePerJob = totalTime / time.Duration(final_completed)
	}

	return LoadTestResult{
		TotalJobs:       numJobs,
		CompletedJobs:   final_completed,
		FailedJobs:      final_failed,
		ProcessingTime:  totalTime,
		JobsPerSecond:   jobsPerSecond,
		AvgTimePerJob:   avgTimePerJob,
		BrowserPoolSize: poolSize,
		WorkerThreads:   maxWorkers,
	}
}

func createMockDelivery(messageId string, body []byte, wg *sync.WaitGroup, mu *sync.Mutex, success, failed *int) amqp.Delivery {
	return amqp.Delivery{
		MessageId: messageId,
		Body:      body,
		Headers:   make(amqp.Table), // Empty headers for clean test
		Acknowledger: &mockAcknowledger{
			wg:      wg,
			mu:      mu,
			success: success,
			failed:  failed,
			id:      messageId,
		},
	}
}

type mockAcknowledger struct {
	wg      *sync.WaitGroup
	mu      *sync.Mutex
	success *int
	failed  *int
	id      string
}

func (m *mockAcknowledger) Ack(tag uint64, multiple bool) error {
	m.mu.Lock()
	*m.success++
	fmt.Printf("✅ Job %s completed successfully\n", m.id)
	m.mu.Unlock()
	m.wg.Done()
	return nil
}

func (m *mockAcknowledger) Nack(tag uint64, multiple, requeue bool) error {
	m.mu.Lock()
	*m.failed++
	if requeue {
		fmt.Printf("🔄 Job %s failed, requeued\n", m.id)
	} else {
		fmt.Printf("❌ Job %s failed permanently\n", m.id)
	}
	m.mu.Unlock()
	m.wg.Done()
	return nil
}

func (m *mockAcknowledger) Reject(tag uint64, requeue bool) error {
	return m.Nack(tag, false, requeue)
}

func printLoadTestResult(result LoadTestResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 YOUR EXACT IMPLEMENTATION LOAD TEST RESULTS")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("🔧 Configuration (YOUR settings):\n")
	fmt.Printf("   Browser Pool Size: %d\n", result.BrowserPoolSize)
	fmt.Printf("   Worker Threads: %d\n", result.WorkerThreads)

	fmt.Printf("\n📊 Results:\n")
	fmt.Printf("   Total Jobs: %d\n", result.TotalJobs)
	fmt.Printf("   Completed: %d\n", result.CompletedJobs)
	fmt.Printf("   Failed: %d\n", result.FailedJobs)
	fmt.Printf("   Success Rate: %.1f%%\n", float64(result.CompletedJobs)/float64(result.TotalJobs)*100)

	fmt.Printf("\n⏱️ Performance:\n")
	fmt.Printf("   Total Processing Time: %v\n", result.ProcessingTime.Round(time.Millisecond))
	fmt.Printf("   Jobs per Second: %.2f\n", result.JobsPerSecond)
	fmt.Printf("   Average Time per Job: %v\n", result.AvgTimePerJob.Round(time.Millisecond))

	// Calculate theoretical vs actual performance
	theoreticalMaxJPS := float64(result.BrowserPoolSize) / 4.0 // Assuming 4 seconds per job
	efficiency := (result.JobsPerSecond / theoreticalMaxJPS) * 100

	fmt.Printf("\n📈 Analysis:\n")
	fmt.Printf("   Theoretical Max: %.2f jobs/sec\n", theoreticalMaxJPS)
	fmt.Printf("   Actual Performance: %.2f jobs/sec\n", result.JobsPerSecond)
	fmt.Printf("   Browser Pool Efficiency: %.1f%%\n", efficiency)

	if result.WorkerThreads > result.BrowserPoolSize {
		fmt.Printf("   ⚠️ More workers (%d) than browsers (%d) - workers waiting for browsers\n",
			result.WorkerThreads, result.BrowserPoolSize)
	}

	fmt.Printf("\n💡 Recommendations:\n")
	if efficiency < 50 {
		fmt.Printf("   🔴 Low efficiency - check for browser deadlocks or failures\n")
	} else if efficiency < 80 {
		fmt.Printf("   🟡 Moderate efficiency - room for improvement\n")
	} else {
		fmt.Printf("   🟢 Good efficiency - browser pool is working well\n")
	}

	// Over-engineering assessment
	fmt.Printf("\n🤔 Over-engineering Assessment:\n")
	if result.JobsPerSecond > 0.75 && efficiency > 70 {
		fmt.Printf("   ✅ Browser pool approach is justified\n")
		fmt.Printf("   - Good throughput with YOUR implementation\n")
	} else if result.JobsPerSecond < 0.3 {
		fmt.Printf("   ❓ Consider simpler approach\n")
		fmt.Printf("   - Single browser might be sufficient for your load\n")
	} else {
		fmt.Printf("   ⚖️ Mixed results - depends on your expected load\n")
	}
}
