package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DanielPopoola/cencori-go"
)

func main() {
	client, err := cencori.NewClient(
		cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get metrics for different periods
	periods := []string{"24h", "7d", "30d"}

	for _, period := range periods {
		fmt.Printf("📊 Metrics for last %s\n", period)

		metrics, err := client.Metrics.Get(context.Background(), period)
		if err != nil {
			log.Printf("Failed to get %s metrics: %v", period, err)
			continue
		}

		printMetrics(metrics)
	}
}

func printMetrics(m *cencori.MetricsResponse) {
	// Request metrics
	fmt.Printf("\n Requests\n")
	fmt.Printf("  Total:        %d\n", m.Requests.Total)
	fmt.Printf("  Successful:   %d (%.1f%%)\n",
		m.Requests.Success, m.Requests.SuccessRate)
	fmt.Printf("  Errors:       %d\n", m.Requests.Error)
	fmt.Printf("  Filtered:     %d\n", m.Requests.Filtered)

	// Cost metrics
	fmt.Printf("\n Costs\n")
	fmt.Printf("  Total:        $%.4f\n", m.Cost.TotalUSD)
	fmt.Printf("  Per Request:  $%.6f\n", m.Cost.AveragePerRequestUSD)

	// Token usage
	fmt.Printf("\n Tokens\n")
	fmt.Printf("  Prompt:       %s\n", formatNumber(m.Tokens.Prompt))
	fmt.Printf("  Completion:   %s\n", formatNumber(m.Tokens.Completion))
	fmt.Printf("  Total:        %s\n", formatNumber(m.Tokens.Total))

	// Latency stats
	fmt.Printf("\n Latency (ms)\n")
	fmt.Printf("  Average:      %d\n", m.Latency.AvgMS)
	fmt.Printf("  P50:          %d\n", m.Latency.P50MS)
	fmt.Printf("  P90:          %d\n", m.Latency.P90MS)
	fmt.Printf("  P99:          %d\n", m.Latency.P99MS)

	// Provider breakdown
	if len(m.Providers) > 0 {
		fmt.Printf("\n By Provider\n")
		for provider, stats := range m.Providers {
			fmt.Printf("  %-12s %d requests, $%.4f\n",
				provider, stats.Requests, stats.CostUSD)
		}
	}

	// Model breakdown
	if len(m.Models) > 0 {
		fmt.Printf("\n By Model\n")
		for model, stats := range m.Models {
			fmt.Printf("  %-20s %d requests, $%.4f\n",
				model, stats.Requests, stats.CostUSD)
		}
	}

	// Cost alerts
	if m.Cost.TotalUSD > 100 {
		fmt.Printf("\n WARNING: Spend exceeds $100\n")
	}
	if m.Requests.SuccessRate < 95 {
		fmt.Printf("\n WARNING: Success rate below 95%%\n")
	}
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}
