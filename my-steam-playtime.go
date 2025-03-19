package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

// SteamVars holds the Steam API key and Steam ID
type SteamVars struct {
	APIKey  string
	SteamID string
}

// Game represents a Steam game with its playtime
type Game struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

// GameOutput represents the output format for a game
type GameOutput struct {
	Rank              int     `json:"rank"`
	GameName          string  `json:"game_name"`
	TotalPlaytimeHours float64 `json:"total_playtime_hours"`
}

// Response from Steam API
type Response struct {
	Response struct {
		Games []Game `json:"games"`
	} `json:"response"`
}

// readSteamVars reads the Steam API key and Steam ID from a file
func readSteamVars(filePath string) (SteamVars, error) {
	var vars SteamVars
	
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return vars, fmt.Errorf("error: %s not found", filePath)
	}
	
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return vars, fmt.Errorf("error: invalid format in %s. Expected 'KEY=VALUE' format", filePath)
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		if key == "STEAM_API_KEY" {
			vars.APIKey = value
		} else if key == "STEAM_ID" {
			vars.SteamID = value
		}
	}
	
	return vars, nil
}

// fetchSteamGames fetches games data from Steam API
func fetchSteamGames(apiKey, steamID string) ([]Game, error) {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", apiKey, steamID)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Steam API: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}
	
	return response.Response.Games, nil
}

// generateMarkdownTable generates a markdown table for the top games sorted by playtime
func generateMarkdownTable(games []Game, numResults int) string {
	if len(games) > numResults {
		games = games[:numResults]
	}
	
	var sb strings.Builder
	sb.WriteString("| Rank | Game Name | Total Playtime (Hours) |\n")
	sb.WriteString("|------|-----------|-------------------------|\n")
	
	for i, game := range games {
		playtimeHours := float64(game.PlaytimeForever) / 60.0
		sb.WriteString(fmt.Sprintf("| %d | %s | %.1f |\n", i+1, game.Name, playtimeHours))
	}
	
	return sb.String()
}

// generateHTMLTable generates an HTML table for the top games sorted by playtime
func generateHTMLTable(games []Game, numResults int) string {
	if len(games) > numResults {
		games = games[:numResults]
	}
	
	var sb strings.Builder
	sb.WriteString("<table>\n")
	sb.WriteString("  <tr>\n")
	sb.WriteString("    <th>Rank</th>\n")
	sb.WriteString("    <th>Game Name</th>\n")
	sb.WriteString("    <th>Total Playtime (Hours)</th>\n")
	sb.WriteString("  </tr>\n")
	
	for i, game := range games {
		playtimeHours := float64(game.PlaytimeForever) / 60.0
		sb.WriteString("  <tr>\n")
		sb.WriteString(fmt.Sprintf("    <td>%d</td>\n", i+1))
		sb.WriteString(fmt.Sprintf("    <td>%s</td>\n", game.Name))
		sb.WriteString(fmt.Sprintf("    <td>%.1f</td>\n", playtimeHours))
		sb.WriteString("  </tr>\n")
	}
	
	sb.WriteString("</table>")
	return sb.String()
}

// generateJSONOutput generates a JSON output for the top games sorted by playtime
func generateJSONOutput(games []Game, numResults int) (string, error) {
	if len(games) > numResults {
		games = games[:numResults]
	}
	
	output := make([]GameOutput, len(games))
	for i, game := range games {
		output[i] = GameOutput{
			Rank:              i + 1,
			GameName:          game.Name,
			TotalPlaytimeHours: float64(game.PlaytimeForever) / 60.0,
		}
	}
	
	jsonData, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON: %v", err)
	}
	
	return string(jsonData), nil
}

func main() {
	// Parse command-line arguments
	numResults := flag.Int("n", 15, "Number of top games to display (default: 15)")
	format := flag.String("format", "markdown", "Output format for the table (markdown, html, or json)")
	flag.Parse()
	
	// Read Steam API key and Steam ID from steam_vars.txt
	steamVars, err := readSteamVars("steam_vars.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	if steamVars.APIKey == "" || steamVars.SteamID == "" {
		fmt.Println("Error: STEAM_API_KEY or STEAM_ID is missing in steam_vars.txt.")
		os.Exit(1)
	}
	
	// Fetch games data from Steam API
	games, err := fetchSteamGames(steamVars.APIKey, steamVars.SteamID)
	if err != nil {
		fmt.Printf("Error fetching games: %v\n", err)
		os.Exit(1)
	}
	
	if len(games) == 0 {
		fmt.Println("No games found or failed to retrieve games.")
		os.Exit(1)
	}
	
	// Sort games by playtime (descending)
	sort.Slice(games, func(i, j int) bool {
		return games[i].PlaytimeForever > games[j].PlaytimeForever
	})
	
	// Generate and print the table in the specified format
	var output string
	switch *format {
	case "markdown":
		output = generateMarkdownTable(games, *numResults)
	case "html":
		output = generateHTMLTable(games, *numResults)
	case "json":
		output, err = generateJSONOutput(games, *numResults)
		if err != nil {
			fmt.Printf("Error generating JSON: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unsupported format: %s\n", *format)
		os.Exit(1)
	}
	
	fmt.Println(output)
}
