### A simple concurrent HTML crawler that scrapes URLs from a starting page and follows links recursively.

Features:
- Starts from a given URL and scrapes all links.
- Uses a worker pool for concurrent scraping to speed up the process.
- Avoids revisiting URLs using a sync.Map.
- Supports a timeout for limiting how long the crawler runs.