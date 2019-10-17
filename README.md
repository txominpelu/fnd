
# FND


**fnd** is heavily inspired by [fzf](https://github.com/junegunn/fzf) but it parses lines as json documents. This allows to search entries by field, output the value of a given field and many other field specific manipulations.

# Basic usage

Quickly search through a list of entries parsed from stdin.

1. Get a list of things.
2. Search through them, filter them, pick one.
3. Do something with it

See [Full Example](#full-example)

# Features

- Contains parsers to convert your input into documents. `line_format`.

    - JSON

    ```bash
    curl https://api.exchangerate-api.com/v4/latest/EUR | \
        jq -c '.rates | to_entries[] | {"Currency": .key, "Rate": .value }' | \
            fnd --line_format json
    ```

    ![Search currency rate](https://github.com/txominpelu/fnd/raw/master/doc/images/currency_json_example.jpg)

    &nbsp;
    - Tabular

    ```bash
    ps aux | fnd --line_format tabular
    ```

    ![Choose a process with fnd](https://github.com/txominpelu/fnd/raw/master/doc/images/tabular_ps_example.jpg)


    &nbsp;
    - Plain

- Customized command output:

    - Choose wich column to output: `--output_column`

    ```bash
    # Kill the selected process
    kill -9 $(ps aux | fnd --line_format tabular --output_column 'PID')
    # fnd will output the PID column
    ```

    - Choose output format with golang templates: `--output_template`

    ```bash
    echo $(ps aux | fnd --line_format tabular --output_template '{{.PID}}-{{.USER}}')
    # fnd will output PID-USER values
    ```

# Examples

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


# Full Example

1. Get a list of running processes.


```bash
ps aux
```

2. Pass them to **fnd** to allow searching through them.


```bash
ps aux | fnd
```

3. Parse input to obtain processes' fields

```bash
ps aux | fnd --line_format tabular
```

4. Do something with the output (E.g print the PID)


```bash
echo $(ps aux | fnd --line_format tabular --output_column)
```


![](https://github.com/txominpelu/fnd/blob/master/doc/videos/fnd-ps-aux.gif)


See other examples at [commands/](commands/)

# [License](LICENSE)

The MIT License (MIT)

Copyright (c) 2019 Íñigo Mediavilla Saiz

