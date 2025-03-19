import json  # Add import for JSON

# ...existing code...

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