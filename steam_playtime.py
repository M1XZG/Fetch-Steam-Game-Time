#!/usr/bin/env python3

import requests
import argparse

# Function to read Steam API key and Steam ID from a file
def load_steam_vars(filename="steam_vars.txt"):
    steam_vars = {}
    try:
        with open(filename, "r") as file:
            for line in file:
                key, value = line.strip().split("=")
                steam_vars[key] = value
    except FileNotFoundError:
        print(f"Error: '{filename}' not found. Please create it with your Steam API key and Steam ID.")
        exit(1)
    except ValueError:
        print(f"Error: Invalid format in '{filename}'. Ensure it's formatted as KEY=VALUE.")
        exit(1)
    
    return steam_vars.get("STEAM_API_KEY"), steam_vars.get("STEAM_ID")

# Load API key and Steam ID
STEAM_API_KEY, STEAM_ID = load_steam_vars()

def get_playtime(steam_id, app_id, api_key):
    url = "http://api.steampowered.com/IPlayerService/GetOwnedGames/v1/"
    params = {
        "key": api_key,
        "steamid": steam_id,
        "include_played_free_games": True,
        "format": "json"
    }
    
    response = requests.get(url, params=params)
    
    if response.status_code != 200:
        print("Error: Failed to fetch data from Steam API.")
        return None
    
    data = response.json()
    
    if "response" in data and "games" in data["response"]:
        games = data["response"]["games"]
        for game in games:
            if game["appid"] == int(app_id):
                playtime_minutes = game["playtime_forever"]
                playtime_hours = round(playtime_minutes / 60, 2)
                return playtime_hours
    
    print("Game not found in the user's library.")
    return None

def main():
    parser = argparse.ArgumentParser(description="Get total playtime for a specific Steam game.")
    parser.add_argument("app_id", type=int, help="The Steam App ID of the game")
    
    args = parser.parse_args()
    
    playtime = get_playtime(STEAM_ID, args.app_id, STEAM_API_KEY)
    if playtime is not None:
        print(f"Total playtime for the game (App ID {args.app_id}): {playtime} hours")

if __name__ == "__main__":
    main()