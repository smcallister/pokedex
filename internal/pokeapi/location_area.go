package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smcallister/pokedex/internal/pokecache"
)

type LocationAreaPage struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string     `json:"name"`
		URL  string     `json:"url"`
	}                   `json:"results"`
}

const LocationAreaURL = "https://pokeapi.co/api/v2/location-area"

var cache = pokecache.NewCache(time.Duration(10 * time.Minute))

func GetLocationAreas(url string) (*LocationAreaPage, error) {
	// Check the cache first.
	if cached, found := cache.Get(url); found {
		var page LocationAreaPage
		if err := json.Unmarshal(cached, &page); err != nil {
			return nil, err
		}
	
		return &page, nil
	}

	// Query for the page of location areas.
	resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch location areas: %s", resp.Status)
	}

	// Convert the result to a page.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body for location areas")
	}

	// Decode the JSON response into a LocationAreaPage struct.
	var page LocationAreaPage
	err = json.Unmarshal(body, &page)
    if err != nil {
        return nil, err
    }

	// Cache the result and return the page.
	cache.Add(url, body)
	return &page, nil
}