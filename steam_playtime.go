package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// SteamVars holds the Steam API key and Steam ID
type SteamVars struct {
	APIKey  string
	SteamID string
}

// Game represents a Steam game
type Game struct {
	AppID           int `json:"appid"`
	PlaytimeForever int `json:"playtime_forever"`
}

// LoadSteamVars loads Steam API key and Steam ID from a file
func LoadSteamVars(filename string) (SteamVars, error) {
	file, err := os.Open(filename)
	if err != nil {
		return SteamVars{}, fmt.Errorf("configuration file '%s' not found", filename)
	}
	defer file.Close()

	vars := SteamVars{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			if key == "STEAM_API_KEY" {
				vars.APIKey = value
			} else if key == "STEAM_ID" {
				vars.SteamID = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return SteamVars{}, err
	}

	return vars, nil
}

// GetPlaytime retrieves the playtime for a specific game by App ID
func GetPlaytime(steamID string, appID int, apiKey string) (float64, error) {
	url := "https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/"
	params := fmt.Sprintf("?key=%s&steamid=%s&include_played_free_games=true&format=json", apiKey, steamID)
	resp, err := http.Get(url + params)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch data from Steam API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch data from Steam API: %s", resp.Status)
	}

	var result struct {
		Response struct {
			Games []Game `json:"games"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	for _, game := range result.Response.Games {
		if game.AppID == appID {
			playtimeHours := float64(game.PlaytimeForever) / 60.0
			return playtimeHours, nil
		}
	}

	return 0, fmt.Errorf("game not found in the user's library")
}

func main() {
	// Parse command-line arguments
	appID := flag.Int("app_id", 0, "The Steam App ID of the game")
	flag.Parse()

	if *appID == 0 {
		fmt.Println("Error: You must provide a valid App ID using the -app_id flag.")
		return
	}

	// Load Steam API key and Steam ID
	steamVars, err := LoadSteamVars("steam_vars.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if steamVars.APIKey == "" || steamVars.SteamID == "" {
		fmt.Println("Error: Steam API Key or Steam ID is missing in steam_vars.txt")
		return
	}

	// Get playtime for the specified game
	playtime, err := GetPlaytime(steamVars.SteamID, *appID, steamVars.APIKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Total playtime for the game (App ID %d): %.2f hours\n", *appID, playtime)
}
