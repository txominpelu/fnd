
# Run it with
$(pipenv --py) run.py  | fzf --inline-info --preview 'echo {} | jq -r ".description" '
