package main

import (
	"fmt"
	"os"

	"github.com/mickamy/contest/example/internal/api"
)

func main() {
	s := api.NewServer()
	if err := s.ListenAndServe(); err != nil {
		fmt.Println("failed to start server:", err)
		os.Exit(1)
	}
}
