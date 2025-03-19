package main

import (
	coding/json"
	ag"
	t"
	ml"
	/ioutil"
	g"
	t/http"
	ring
	
type Game struct {
    Name            string `json:"name"`
    PlaytimeForever int    `json:"playtime_forever"`
}

type SteamVars struct {
    APIKey  string
    SteamID string
}

func readSteamVars(filePath string) SteamVars {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        log.Fatalf("Error: %s not found.\n", filePath)
    }

    lines := strings.Split(string(data), "\n")
    vars := SteamVars{}

    for _, line := range lines {
        if line == "" {
            continue
        }
        parts := strings.Split(line, "=")
        if len(parts) != 2 {
            log.Fatalf("Error: Invalid format in %s. Expected 'KEY=VALUE' format.\n", filePath)
        }
        key, value := parts[0], parts[1]
        if key == "STEAM_API_KEY" {
            vars.APIKey = value
        } else if key == "STEAM_ID" {
            vars.SteamID = value
        }
    }

    if vars.APIKey == "" || vars.SteamID == "" {
        log.Fatalf("Error: STEAM_API_KEY or STEAM_ID is missing in %s.\n", filePath)
    }

    return vars
}

func fetchSteamGames(apiKey, steamID string) []Game {
    url := "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatalf("Error creating request: %v\n", err)
    }

    query := req.URL.Query()
    query.Add("key", apiKey)
    query.Add("steamid", steamID)
    query.Add("include_appinfo", "true")
    query.Add("format", "json")
    req.URL.RawQuery = query.Encode()

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatalf("Error fetching data from Steam API: %v\n", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("Error: Failed to fetch data from Steam API. Status code: %d\n", resp.StatusCode)
    }

    var result struct {
        Response struct {
            Games []Game `json:"games"`
        } `json:"response"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Fatalf("Error decoding response: %v\n", err)
    }

    return result.Response.Games
}

func generateMarkdownTable(games []Game, numResults int) string {
    sort.Slice(games, func(i, j int) bool {
        return games[i].PlaytimeForever > games[j].PlaytimeForever
    })

    if len(games) > numResults {
        games = games[:numResults]
    }

    table := "| Rank | Game Name | Total Playtime (Hours) |\n"
    table += "|------|-----------|-------------------------|\n"

    for i, game := range games {
        playtimeHours := float64(game.PlaytimeForever) / 60
        table += fmt.Sprintf("| %d | %s | %.1f |\n", i+1, game.Name, playtimeHours)
    }

    return table
}

func generateHTMLTable(games []Game, numResults int) string {
    sort.Slice(games, func(i, j int) bool {
        return games[i].PlaytimeForever > games[j].PlaytimeForever
    })

    if len(games) > numResults {
        games = games[:numResults]
    }

    table := "<table>\n"
    table += "  <tr><th>Rank</th><th>Game Name</th><th>Total Playtime (Hours)</th></tr>\n"

    for i, game := range games {
        playtimeHours := float64(game.PlaytimeForever) / 60
        table += fmt.Sprintf("  <tr><td>%d</td><td>%s</td><td>%.1f</td></tr>\n", i+1, html.EscapeString(game.Name), playtimeHours)
    }

    table += "</table>"
    return table
}

func generateJSONOutput(games []Game, numResults int) string {
    sort.Slice(games, func(i, j int) bool {
        return games[i].PlaytimeForever > games[j].PlaytimeForever
    })

    if len(games) > numResults {
        games = games[:numResults]
    }

    output := []map[string]interface{}{}
    for i, game := range games {
        playtimeHours := float64(game.PlaytimeForever) / 60
        output = append(output, map[string]interface{}{
            "rank":                i + 1,
            "game_name":           game.Name,
            "total_playtime_hours": playtimeHours,
        })
    }

    jsonData, err := json.MarshalIndent(output, "", "    ")
    if err != nil {
        log.Fatalf("Error generating JSON: %v\n", err)
    }

    return string(jsonData)
}

func main() {
    numResults := flag.Int("n", 15, "Number of top games to display (default: 15)")
    format := flag.String("format", "markdown", "Output format for the table (markdown, html, json)")
    flag.Parse()

    steamVars := readSteamVars("steam_vars.txt")
    games := fetchSteamGames(steamVars.APIKey, steamVars.SteamID)

    if len(games) == 0 {
        log.Fatalf("No games found or failed to retrieve games.\n")
    }

    var output string
    switch *format {
    case "markdown":
        output = generateMarkdownTable(games, *numResults)
    case "html":
        output = generateHTMLTable(games, *numResults)
    case "json":
        output = generateJSONOutput(games, *numResults)
    default:
        log.Fatalf("Invalid format: %s. Valid options are markdown, html, json.\n", *format)
    }

    fmt.Println(output)
}