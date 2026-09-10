package internal

import "github.com/mrizwan/pokedexcli/internal/pokeapi"

type Pokedex map[string]pokeapi.Pokemon

func NewPokedex() *Pokedex {
	pokedex := make(Pokedex)
	return &pokedex
}

func (p Pokedex) Add(pokemon pokeapi.Pokemon) {
	p[pokemon.Name] = pokemon
}

func (p Pokedex) Get(name string) (pokeapi.Pokemon, bool) {
	pokemon, ok := p[name]
	return pokemon, ok
}
