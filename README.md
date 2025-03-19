# Steam Playtime Query Script

This script fetches the total playtime for a specific game from the Steam API.

## Features
- Reads `STEAM_API_KEY` and `STEAM_ID` from `steam_vars.txt`
- Fetches total playtime for a game using its **App ID**
- Supports direct execution (`chmod +x steam_playtime.py`)
- Uses **command-line arguments** for flexibility

---

## Prerequisites
1. **Get a Steam API Key** from [Steam Developer Portal](https://steamcommunity.com/dev/apikey).
2. **Find your Steam ID** (64-bit format) from [Steam ID Finder](https://steamid.io/).
3. **Find the game's App ID** from [SteamDB](https://steamdb.info/).
4. **Create a `steam_vars.txt` file** in the same directory as the script:

   ```txt
   STEAM_API_KEY=your_actual_api_key
   STEAM_ID=your_actual_steam_id
   ```

---

## Installation & Setup
1. Clone the repository or create a new script file:
   ```sh
   git clone https://github.com/yourusername/steam-playtime-query.git
   cd steam-playtime-query
   ```
2. Make the script executable (Linux/macOS):
   ```sh
   chmod +x steam_playtime.py
   ```

---

## Usage
Run the script with the game's **App ID** as a command-line argument:

### Method 1: Direct Execution
```sh
./steam_playtime.py 730
```
*(Replace `730` with the App ID of the game.)*

### Method 2: Using Python
```sh
python3 steam_playtime.py 730
```

---

## Example Output
```
Total playtime for the game (App ID 730): 125.5 hours
```

---

## Script (`steam_playtime.py`)
```python
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
```

---

## License
MIT License (or add your preferred license here).

## Author
[Your Name](https://github.com/yourusername/)
