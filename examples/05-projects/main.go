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

	orgSlug := "my-org" // Your organization slug

	// List all projects
	fmt.Println("Listing projects...")
	projects, err := client.Projects.List(context.Background(), orgSlug)
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	fmt.Printf("Found %d projects:\n", len(projects))
	for _, p := range projects {
		fmt.Printf("  - %s (%s)\n", p.Name, p.Slug)
	}

	// Create new project
	fmt.Println("Creating new project...")
	newProject, err := client.Projects.Create(context.Background(), orgSlug, cencori.CreateProjectParams{
		Name:        "AI Assistant",
		Description: "Customer support chatbot",
		Visibility:  "public",
	})
	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}

	fmt.Printf("Created: %s (ID: %s)\n", newProject.Name, newProject.ID)

	// Get project details
	fmt.Println("Fetching project details...")
	project, err := client.Projects.Get(context.Background(), orgSlug, newProject.Slug)
	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}

	fmt.Printf("Project: %s\n", project.Name)
	fmt.Printf("Status: %s\n", project.Status)
	fmt.Printf("Created: %s\n", project.CreatedAt)
	if project.Stats != nil {
		fmt.Printf("Requests: %d\n", project.Stats.TotalRequests)
		fmt.Printf("Cost: $%.2f\n", project.Stats.TotalCostUSD)
	}

	// Clean up - delete project
	// Uncomment to actually delete:
	// fmt.Println("Deleting project...")
	// if err := client.Projects.Delete(context.Background(), orgSlug, newProject.Slug); err != nil {
	//     log.Fatalf("Failed to delete project: %v", err)
	// }
	// fmt.Println("Project deleted")
}
