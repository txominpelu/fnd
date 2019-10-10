
TODO:

- Parse whole string and allow jq path selection for json instead of parsing line by line
- Show header other than $ in plain text format
- Try the trie to avoid having to match exactly on word
- When tokenizing don't split by dot, just stem by it
- Log to a ~/.fnd logs file

Features:

- Display headers at the top (requires tabular printing for alignment)
- Enter - Returns currently selected item
- Streaming entries - Show while continue reading stdin
- Select entry with up - down
