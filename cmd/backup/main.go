package main

import (
    "fmt"
    "log"
    "github.com/joho/godotenv"
    "github.com/lamdaloop/kubedock/internal/k8s"
)

func main() {
    // Load environment variables from .env
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  .env file not found, proceeding with system env")
    }

    fmt.Println("🚀 Starting KubeDock Backup...")

    _, discoClient, dynClient := k8s.CreateClient()

    resources, err := k8s.DiscoverResources(discoClient)
    if err != nil {
        log.Fatalf("❌ Failed to discover resources: %v", err)
    }

    err = k8s.FetchAndDumpResources(dynClient, resources, "./backups")
    if err != nil {
        log.Fatalf("❌ Backup failed: %v", err)
    }

    fmt.Println("✅ Backup completed successfully.")
}
