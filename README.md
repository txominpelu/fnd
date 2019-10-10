
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

If nothing is passed in stdin:

1) If folder inside a git repo, the list contains the result of *git ls-files*.
2) If it's not inside a git repo, the list contains a recursive walk of all the files inside the current dir.

Examples:

- Open file with vi:

    ```
    $> vi $(fnd)
    ```

- Pass tabular format (pass a table with a header and separated by spaces )

    ```
    $> ps | fnd --line_format tabular
    ```

- You can choose which column will be the output of the command. E.g this is how to kill the process that is chosen in fnd.

    ```
    $> ps | kill -9 $(fnd --line_format tabular --output_column PID)
    ```

# Use case 

```
fnd --see '.description' --query 'description:hello'
```
