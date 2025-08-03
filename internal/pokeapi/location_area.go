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

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string 				`json:"name"`
			URL  string 				`json:"url"`
		} 								`json:"encounter_method"`
		VersionDetails []struct {
			Rate    int 				`json:"rate"`
			Version struct {
				Name string 			`json:"name"`
				URL  string 			`json:"url"`
			} 							`json:"version"`
		} 								`json:"version_details"`
	} 									`json:"encounter_method_rates"`
	GameIndex int 						`json:"game_index"`
	ID        int 						`json:"id"`
	Location  struct {
		Name string 					`json:"name"`
		URL  string 					`json:"url"`
	} 									`json:"location"`
	Name  string 						`json:"name"`
	Names []struct {
		Language struct {
			Name string 				`json:"name"`
			URL  string 				`json:"url"`
		} 								`json:"language"`
		Name string 					`json:"name"`
	} 									`json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string 				`json:"name"`
			URL  string 				`json:"url"`
		} 								`json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   	`json:"chance"`
				ConditionValues []any 	`json:"condition_values"`
				MaxLevel        int   	`json:"max_level"`
				Method          struct {
					Name string 		`json:"name"`
					URL  string 		`json:"url"`
				} 						`json:"method"`
				MinLevel int 			`json:"min_level"`
			} 							`json:"encounter_details"`
			MaxChance int 				`json:"max_chance"`
			Version   struct {
				Name string 			`json:"name"`
				URL  string 			`json:"url"`
			} 							`json:"version"`
		} 								`json:"version_details"`
	} 									`json:"pokemon_encounters"`
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

func GetLocationArea(name string) (*LocationArea, error) {
	// Construct the URL for the specific location area.
	url := fmt.Sprintf("%s/%s", LocationAreaURL, name)

	// Check the cache first.
	if cached, found := cache.Get(url); found {
		var area LocationArea
		if err := json.Unmarshal(cached, &area); err != nil {
			return nil, err
		}
	
		return &area, nil
	}

	// Query for the specific location area.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch location area: %s", resp.Status)
	}

	// Convert the result to a LocationArea struct.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("Empty response body for location area")
	}

	// Decode the JSON response into a LocationArea struct.
	var area LocationArea
	err = json.Unmarshal(body, &area)
	if err != nil {
		return nil, err
	}

	// Cache the result and return the area.
	cache.Add(url, body)
	return &area, nil
}