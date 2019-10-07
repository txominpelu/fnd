
# fnd

Quickly search through a list, choose elements. 

fnd is heavily inspired by [fzf](https://github.com/junegunn/fzf) but with additional features:

- Contains parsers to give structure to your lines.
    - json - WIP
    - tabular format - WIP

- Customized input and output with jq.
    - Choose how your input will look like. - WIP

        ```
        # parse lines as a tabular with header and choose which columns to display 
        $> ps aux | fnd --parser tabular --header --display-columns "USER,PID,%CPU,%MEM,COMMAND" 
        ```
    - Choose what the output will look like - WIP
        ```
        # Kill the selected process
        $> kill -9 $(ps aux | fnd --parser tabular --header --output '.PID')
        ```
- Pick multiple options - WIP


## Basic usage

1) Pipe a list to fnd  2) search elements 3) pick your favourite or exit. 

Examples:

- Open file with vi:

    ```
    $> 
    ```

# Run it with
$(pipenv --py) run.py  | fzf --inline-info --preview 'echo {} | jq -r ".description" '

# Use case 

./rjobs --see '.description' --query 'description:hello'
