
TODO:

- Don't kill the world when it fails (e.g if rg is not installed fnd-rg-edit kills the current iterm tab)
- Pick multiple entries (copy multiple files to another folder)
- Sort by column (asc, desc) - Interactive 
- Sort by column (asc, desc) - CLI 
- SQL like queries
- Tokenize queries main.go should search for query and go (or define expectations for search altogether)
- When tokenizing don't split by dot, just stem by it
- Scroll up and down through results
- Show header other than $ in plain text format
- Index by char position for fuzzy search
- Try the trie to avoid having to match exactly on word

Features:

- Display headers at the top (requires tabular printing for alignment)
- Enter - Returns currently selected item
- Streaming entries - Show while continue reading stdin
- Select entry with up - down
- Log errors to stderr or specified log file 
