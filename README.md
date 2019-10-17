
```
   __           _ 
  / ■|         |■|
 |■|_ _ __   __|■|
 |■ ■|■'■ \ /■■`■|
 |■| |■| |■|■(_|■|
 |■| |■| |■|\■■,■|

```


Quickly search through a list of documents parsed from stdin. 

fnd is heavily inspired by [fzf](https://github.com/junegunn/fzf) but with additional features:


# Aditional features

- Contains parsers to convert your input into documents. `line_format`.

    - json

    ```bash
    curl https://api.exchangerate-api.com/v4/latest/EUR | \
        jq -c '.rates | to_entries[] | {"Currency": .key, "Rate": .value }' | \
            fnd --line_format json
    ```

    <kbd>
        <img src="https://github.com/txominpelu/fnd/raw/master/doc/images/currency_json_example.png" alt="Search currency rate">
    </kbd>

    - tabular
    - plain

- Customized command output:

    - Choose what the output will look like
        - Choose wich column to output: `--output_column`

    ```bash
    # Kill the selected process
    kill -9 $(ps aux | fnd --parser tabular --header --output_column 'PID')
    # fnd will output the PID column
    ```

    - Choose output format with golang templates: `--output_template`

    ```bash
    echo $(ps aux | fnd --parser tabular --header --output_template '{{.PID}}-{{.USER}}')
    # fnd will output PID-USER values
    ```


- Pick multiple options - WIP


## Basic usage

1) Pipe a list to fnd  2) search elements 3) pick your favourite or exit. 

If nothing is passed in stdin:

1) If folder inside a git repo, the list contains the result of *git ls-files*.
2) If it's not inside a git repo, the list contains a recursive walk of all the files inside the current dir.

Examples:

- Open file with vi:

    ```bash
    vi $(fnd-fdfind)
    ```

- Pass tabular format (pass a table with a header and separated by spaces )

    ```bash
    ps | fnd --line_format tabular
    ```

- You can choose which column will be the output of the command. E.g this is how to kill the process that is chosen in fnd.

    ```bash
    ps | kill -9 $(fnd --line_format tabular --output_column PID)
    ```

# Troubleshooting

- My command fails, how can I figure out what's happening ?

Error logs are by default logged to stderr. To get the detailed error message you can redirect stderr to a log file or pass --log_file. E.g: 

```bash
$ fnd 2> out.log
```


# Examples

See at [commands/](commands/)

# [License](LICENSE)

The MIT License (MIT)

Copyright (c) 2019 Íñigo Mediavilla Saiz

