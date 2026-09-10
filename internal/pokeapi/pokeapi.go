package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mrizwan/pokedexcli/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2/"
)

// Client -
type Client struct {
	cache      *pokecache.Cache
	httpClient http.Client
}

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	}
}

// NewClient -
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

var next, previous = "", ""

func (c Client) GetLocationAreas(previousAreas bool) ([]string, error) {
	var resp_data []byte

	areasURL := baseURL + "location-area/"
	if previousAreas {
		if previous != "" {
			areasURL = previous
		}
	} else {
		if next != "" {
			areasURL = next
		}
	}

	if data, ok := c.cache.Get(areasURL); ok {
		resp_data = data
	} else {
		resp, err := c.httpClient.Get(areasURL)
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()

		resp_data, err = io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}
		c.cache.Add(areasURL, resp_data)
	}

	var result LocationAreaResponse
	err := json.Unmarshal(resp_data, &result)
	if err != nil {
		return []string{}, err
	}

	next = result.Next
	previous = result.Previous

	var areas []string
	for _, area := range result.Results {
		areas = append(areas, area.Name)
	}
	return areas, nil
}

func (c Client) GetPokemonsInArea(area string) ([]string, error) {
	pokemonURL := baseURL + "location-area/" + area
	var resp_data []byte

	if data, ok := c.cache.Get(pokemonURL); ok {
		resp_data = data
	} else {
		resp, err := c.httpClient.Get(pokemonURL)
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return []string{}, fmt.Errorf("(%s)", resp.Status)
		}

		resp_data, err = io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}
		c.cache.Add(pokemonURL, resp_data)
	}

	var result PokemonResponse

	err := json.Unmarshal(resp_data, &result)
	if err != nil {
		return []string{}, err
	}

	var pokemons []string
	for _, encounter := range result.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon.Name)
	}
	return pokemons, nil
}

func (c Client) GetPokemonDetails(pokemon string) (Pokemon, error) {
	pokemonURL := baseURL + "pokemon/" + pokemon
	var resp_data []byte

	if data, ok := c.cache.Get(pokemonURL); ok {
		resp_data = data
	} else {
		resp, err := c.httpClient.Get(pokemonURL)
		if err != nil {
			return Pokemon{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return Pokemon{}, fmt.Errorf("(%s)", resp.Status)
		}

		resp_data, err = io.ReadAll(resp.Body)
		if err != nil {
			return Pokemon{}, err
		}
		c.cache.Add(pokemonURL, resp_data)
	}

	var result Pokemon

	err := json.Unmarshal(resp_data, &result)
	if err != nil {
		return Pokemon{}, err
	}

	return result, nil
}
