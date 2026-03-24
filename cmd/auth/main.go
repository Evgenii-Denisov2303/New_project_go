package main

import (
	"fmt"
	"new_project_go/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("auth service starting... env=%s port=%s\n", cfg.Env, cfg.HTTPPort)
}
