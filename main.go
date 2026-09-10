package main

import (
	"time"

	"github.com/mrizwan/pokedexcli/internal"
)

func main() {
	cfg := internal.NewConfig(5*time.Second, 5*time.Minute)
	internal.StartRepl(cfg)
}
