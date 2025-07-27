package main

import (
	"fmt"
	"net/http"
)

func main() {
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz/", handleHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleGetFileserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleResetFileserverHits)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server started on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
