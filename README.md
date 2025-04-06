## denv

This is CLI tool to create instant docker contanier for developping.


## install

```
go install github.com/takumi2786/denv@version
```

## how to use

```sh
denv --help
This is a CLI to manipulate instant Docker containers.

Usage:
denv [command]

Available Commands:
completion  Generate the autocompletion script for the specified shell
delete      Delete selected container
help        Help about any command
run         Start Instant Container

Flags:
-f, --file string   path to image_map.json (default "resources/image_map.json")
-h, --help          help for denv

Use "denv [command] --help" for more information about a command.
```

1. Create json file like <a href="./resources/image_map.json">resources/image_map.json</a>

2. Run container and attach.
    ```sh
    denv run -i ubuntu -f ./resources/image_map.json
    ```

3. Delete container
    ```
    denv delete -i ubuntu -f ./resources/image_map.json
    ```
    

