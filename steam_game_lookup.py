#!/usr/bin/env python3

import requests
import os

def load_steam_vars(filename="steam_vars.txt"):
    """Load Steam API key and Steam ID from a file."""
    steam_vars = {}
    if not os.path.exists(filename):
        raise FileNotFoundError(f"Configuration file '{filename}' not found.")

    with open(filename, "r") as file:
        for line in file:
            key, value = line.strip().split("=", 1)
            steam_vars[key] = value

    return steam_vars.get("STEAM_API_KEY"), steam_vars.get("STEAM_ID")

def get_owned_games(api_key, steam_id):
    """Retrieve the list of owned games for a user."""
    url = f"http://api.steampowered.com/IPlayerService/GetOwnedGames/v1/"
    params = {
        "key": api_key,
        "steamid": steam_id,
        "include_appinfo": True
    }
    response = requests.get(url, params=params)

    if response.status_code == 200:
        return response.json().get("response", {}).get("games", [])
    else:
        print("Error fetching game list:", response.text)
        return []

def find_game_id(game_name, games):
    """Find the game ID by game name."""
    for game in games:
        if game_name.lower() in game.get("name", "").lower():
            return game["appid"]
    return None

def find_game_name(game_id, games):
    """Find the game name by game ID."""
    for game in games:
        if game["appid"] == game_id:
            return game["name"]
    return None

def main():
    api_key, steam_id = load_steam_vars()
    if not api_key or not steam_id:
        print("Error: Steam API Key or Steam ID is missing in steam_vars.txt")
        return

    games = get_owned_games(api_key, steam_id)
    if not games:
        print("No games found or error fetching games.")
        return

    search = input("Enter a game name or game ID: ")

    if search.isdigit():
        search_id = int(search)
        result = find_game_name(search_id, games)
        print(f"Game Name: {result}" if result else "Game not found.")
    else:
        result = find_game_id(search, games)
        print(f"Game ID: {result}" if result else "Game not found.")

if __name__ == "__main__":
    main()