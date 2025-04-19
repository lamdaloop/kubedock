package main

import (
    "log"
    "net/http"

    "github.com/lamdaloop/kubedock/backend/api"
    "github.com/rs/cors"
)

func main() {
    api.InitDB()

    router := api.NewRouter()

    // Wrap the router with CORS
    handler := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"}, // frontend origin
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    }).Handler(router)

    log.Println("🚀 Backend running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
