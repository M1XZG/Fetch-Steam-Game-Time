# Steam Playtime Scripts

These scripts will fetch game information using the Steam API.

## steam_game_lookup.py

This will help you find a Game ID or find a game name from a Game ID, examples

```
$ ./steam_game_lookup.py
Enter a game name or game ID: vrchat
Game ID: 438100

$ ./steam_game_lookup.py
Enter a game name or game ID: 438100
Game Name: VRChat
```

## steam_playtime.py

This will take a single argument of the Game ID and give you the total playtime in that game, example

```
$ ./steam_playtime.py 438100
Total playtime for the game (App ID 438100): 8007.92 hours
```

## Features
- Reads `STEAM_API_KEY` and `STEAM_ID` from `steam_vars.txt`
- Fetches steam Game ID by name or steam Game Name by ID
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
   git clone https://github.com/M1XZG/Fetch-Steam-Game-Time.git
   cd Fetch-Steam-Game-Time
   ```
2. Make the script executable (Linux/macOS):
   ```sh
   chmod +x steam_playtime.py
   ```
3. You may also require the `requests` package so I've provided the `requirements.txt` for you
   ```sh
   pip3 install -r requirements.txt
   ```

---

## Usage steam_playtime.py
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
