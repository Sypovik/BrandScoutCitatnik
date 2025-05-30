package main

import (
	"BrandScoutCitatnik/internal/handlers"
	"BrandScoutCitatnik/internal/middleware"
	"BrandScoutCitatnik/internal/repository"
	"BrandScoutCitatnik/internal/router"
	"net/http"
	"os"
)

func main() {
	repo := repository.NewQuoteRepositoryMemory()
	handler := handlers.New(repo)
	router := router.Init(handler, middleware.Logger)

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.ListenAndServe(":"+port, router)
}
