# Steam Playtime Scripts

These scripts will fetch game information using the Steam API.

## steam_game_lookup.py

This script mostly works but I have found it can fail to find some games both by ID or name, I'm putting this down to some weird API stuff and will find a fix for it.

This will help you find a Game ID or find a game name from a Game ID, examples

```
$ ./steam_game_lookup.py
Enter a game name or game ID: vrchat
Game ID: 438100

$ ./steam_game_lookup.py
Enter a game name or game ID: 438100
Game Name: VRChat
```

You can also run this and provide the Game ID or game name on the command line and have it returned

```
$ ./steam_game_lookup.py 438100
Game Name: VRChat

$ ./steam_game_lookup.py vrchat
Game ID: 438100
```

## my-steam-playtime.py

This will collect the hours of game play for every game in your steam library and output a markdown table sorted with most hours at the top. By default this will print the top 15 games, you can override this by using the `-n X` option where X is the number of games to display. If you want to display all games just use something like `99999`. Example output:

` $ ./my-steam-playtime.py -n 10`
| Rank | Game Name | Total Playtime (Hours) |
|------|-----------|-------------------------|
| 1 | VRChat | 8008.4 |
| 2 | OVR Toolkit | 7201.5 |
| 3 | OVR Advanced Settings | 7187.7 |
| 4 | Standable: Full Body Estimation | 4770.9 |
| 5 | fpsVR | 2498.3 |
| 6 | Overwatch® 2 | 592.6 |
| 7 | World of Warships | 425.7 |
| 8 | Rust | 415.1 |
| 9 | Satisfactory | 248.7 |
| 10 | Stardew Valley | 175.9 |

To get additional options run `my-steam-playtime.py -h`

Additional output formats include HTML table and JSON.

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
3. **Rename `steam_vars.TEMPLATE` to `steam_vars.txt`** in the same directory as the script:

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
   chmod +x *.py
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

## Using Go
There is also a version for Go which can be used.
1. Install Go
2. Fill in the `STEAM_API_KEY` and `STEAM_ID` in the `steam_vars.txt` file
3. Build the Go script `go build my-steam-playtime.go`
4. Run the file `./my-steam-playtime`
