package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DanielPopoola/cencori-go"
)

func main() {
	client, err := cencori.NewClient(
		cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	projectID := "proj_123" // Your project ID
	env := "production"

	// 1. List existing keys
	fmt.Println("Current API keys:")
	keys, err := client.APIKeys.List(context.Background(), projectID, env)
	if err != nil {
		log.Fatalf("Failed to list keys: %v", err)
	}

	for _, key := range keys {
		fmt.Printf("  - %s (created: %s, used: %d times)\n",
			key.Name, key.CreatedAt.Format("2006-01-02"), key.UsageCount)
	}

	fmt.Println("Creating new API key...")
	newKey, err := client.APIKeys.Create(context.Background(), projectID, cencori.CreateAPIKeyParams{
		Name:        "Production Key 2025-01",
		Environment: env,
	})
	if err != nil {
		log.Fatalf("Failed to create key: %v", err)
	}

	// IMPORTANT: Save this key immediately - it's only shown once!
	fmt.Printf("\n SAVE THIS KEY NOW:\n%s\n\n", newKey.Key)
	fmt.Println("This is the only time you'll see the full key!")

	// 3. Wait for deployment to use new key
	fmt.Println("\n Waiting for new key to be deployed...")
	fmt.Println("(In real scenario: update env vars, redeploy services)")
	time.Sleep(2 * time.Second)

	// 4. Verify old key usage has stopped
	fmt.Println("\n Checking key usage...")
	for _, oldKey := range keys {
		stats, err := client.APIKeys.GetStats(context.Background(), projectID, oldKey.ID)
		if err != nil {
			log.Printf("Warning: couldn't get stats for %s: %v", oldKey.ID, err)
			continue
		}

		fmt.Printf("Key %s:\n", oldKey.Name)
		fmt.Printf("  Total requests: %d\n", stats.TotalRequests)
		fmt.Printf("  Last used: %s\n", stats.LastUsedAt)

		// If not used recently, safe to revoke
		if time.Since(stats.LastUsedAt) > 48*time.Hour {
			fmt.Printf("  Safe to revoke (not used in 48h)\n")
		}
	}

	// 5. Revoke old keys
	// Uncomment to actually revoke:
	// fmt.Println("\n Revoking old keys...")
	// for _, oldKey := range keys {
	//     if err := client.APIKeys.Revoke(context.Background(), projectID, oldKey.ID); err != nil {
	//         log.Printf("Failed to revoke %s: %v", oldKey.ID, err)
	//         continue
	//     }
	//     fmt.Printf("  ✓ Revoked: %s\n", oldKey.Name)
	// }

	fmt.Println("\n Key rotation complete!")
}
