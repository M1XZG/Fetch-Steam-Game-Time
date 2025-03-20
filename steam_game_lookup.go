package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
)

// SteamVars holds the Steam API key and Steam ID
type SteamVars struct {
    APIKey  string
    SteamID string
}

// Game represents a Steam game
type Game struct {
    AppID int    `json:"appid"`
    Name  string `json:"name"`
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

// GetOwnedGames retrieves the list of owned games for a user
func GetOwnedGames(apiKey, steamID string) ([]Game, error) {
    url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true", apiKey, steamID)
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error fetching game list: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error fetching game list: %s", resp.Status)
    }

    var result struct {
        Response struct {
            Games []Game `json:"games"`
        } `json:"response"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("error decoding response: %v", err)
    }

    return result.Response.Games, nil
}

// FindGameID finds the game ID by game name
func FindGameID(gameName string, games []Game) int {
    for _, game := range games {
        if strings.Contains(strings.ToLower(game.Name), strings.ToLower(gameName)) {
            return game.AppID
        }
    }
    return 0
}

// FindGameName finds the game name by game ID
func FindGameName(gameID int, games []Game) string {
    for _, game := range games {
        if game.AppID == gameID {
            return game.Name
        }
    }
    return ""
}

func main() {
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

    // Fetch owned games
    games, err := GetOwnedGames(steamVars.APIKey, steamVars.SteamID)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    if len(games) == 0 {
        fmt.Println("No games found or error fetching games.")
        return
    }

    // Handle user input
    var search string
    if len(os.Args) > 1 {
        search = strings.Join(os.Args[1:], " ")
    } else {
        fmt.Print("Enter a game name or game ID: ")
        fmt.Scanln(&search)
    }

    if gameID, err := strconv.Atoi(search); err == nil {
        // Search by game ID
        result := FindGameName(gameID, games)
        if result != "" {
            fmt.Printf("Game Name: %s\n", result)
        } else {
            fmt.Println("Game not found.")
        }
    } else {
        // Search by game name
        result := FindGameID(search, games)
        if result != 0 {
            fmt.Printf("Game ID: %d\n", result)
        } else {
            fmt.Println("Game not found.")
        }
    }
}