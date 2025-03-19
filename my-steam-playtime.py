#!/usr/bin/env python3
# filepath: steam_games_report.py

import requests
import argparse
import json  # Add import for JSON

def read_steam_vars(file_path):
    """Reads the Steam API key and Steam ID from a file."""
    steam_vars = {}
    try:
        with open(file_path, 'r') as file:
            for line in file:
                key, value = line.strip().split('=')
                steam_vars[key] = value
    except FileNotFoundError:
        print(f"Error: {file_path} not found.")
        exit(1)
    except ValueError:
        print(f"Error: Invalid format in {file_path}. Expected 'KEY=VALUE' format.")
        exit(1)
    return steam_vars

def fetch_steam_games(api_key, steam_id):
    """Fetches the list of games and playtime from the Steam API."""
    url = f"http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/"
    params = {
        'key': api_key,
        'steamid': steam_id,
        'include_appinfo': True,
        'format': 'json'
    }
    response = requests.get(url, params=params)
    if response.status_code != 200:
        print(f"Error: Failed to fetch data from Steam API. Status code: {response.status_code}")
        exit(1)
    return response.json().get('response', {}).get('games', [])

def generate_markdown_table(games, num_results):
    """Generates a GitHub Markdown table for the top games sorted by playtime."""
    sorted_games = sorted(games, key=lambda x: x['playtime_forever'], reverse=True)[:num_results]
    table = "| Rank | Game Name | Total Playtime (Hours) |\n"
    table += "|------|-----------|-------------------------|\n"
    for rank, game in enumerate(sorted_games, start=1):
        playtime_hours = game['playtime_forever'] / 60  # Convert minutes to hours
        table += f"| {rank} | {game['name']} | {playtime_hours:.1f} |\n"  # Format to 1 decimal place
    return table

def generate_html_table(games, num_results):
    """Generates an HTML table for the top games sorted by playtime."""
    sorted_games = sorted(games, key=lambda x: x['playtime_forever'], reverse=True)[:num_results]
    table = "<table>\n"
    table += "  <tr><th>Rank</th><th>Game Name</th><th>Total Playtime (Hours)</th></tr>\n"
    for rank, game in enumerate(sorted_games, start=1):
        playtime_hours = game['playtime_forever'] / 60  # Convert minutes to hours
        table += f"  <tr><td>{rank}</td><td>{game['name']}</td><td>{playtime_hours:.1f}</td></tr>\n"
    table += "</table>"
    return table

def generate_json_output(games, num_results):
    """Generates a JSON output for the top games sorted by playtime."""
    sorted_games = sorted(games, key=lambda x: x['playtime_forever'], reverse=True)[:num_results]
    output = [
        {
            "rank": rank,
            "game_name": game['name'],
            "total_playtime_hours": round(game['playtime_forever'] / 60, 1)  # Convert minutes to hours
        }
        for rank, game in enumerate(sorted_games, start=1)
    ]
    return json.dumps(output, indent=4)  # Pretty-print JSON with indentation

def main():
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description="Generate a table of your Steam games sorted by playtime.")
    parser.add_argument(
        "-n", "--num-results",
        type=int,
        default=15,
        help="Number of top games to display (default: 15)"
    )
    parser.add_argument(
        "--format",
        choices=["markdown", "html", "json"],  # Add "json" as a valid choice
        default="markdown",
        help="Output format for the table (default: markdown)"
    )
    args = parser.parse_args()

    # Read Steam API key and Steam ID from steam_vars.txt
    steam_vars = read_steam_vars('steam_vars.txt')
    api_key = steam_vars.get('STEAM_API_KEY')
    steam_id = steam_vars.get('STEAM_ID')

    if not api_key or not steam_id:
        print("Error: STEAM_API_KEY or STEAM_ID is missing in steam_vars.txt.")
        exit(1)

    # Fetch games data from Steam API
    games = fetch_steam_games(api_key, steam_id)

    if not games:
        print("No games found or failed to retrieve games.")
        exit(1)

    # Generate and print the table in the specified format
    if args.format == "markdown":
        table = generate_markdown_table(games, args.num_results)
    elif args.format == "html":
        table = generate_html_table(games, args.num_results)
    elif args.format == "json":  # Handle JSON format
        table = generate_json_output(games, args.num_results)
    print(table)

if __name__ == "__main__":
    main()