# tlint

My own dumb linting program, ties together a few common linters with some config files
to do common things in projects.

YOU MUST INSTALL THE LINTERS YOURSELF!
This is just glue to put them together with a config file

Runs:

- goimports
- go fmt
- gometalinter
- hadolint
- shellcheck
- checkmake

recursively under the cwd.
It will glob appropriately for the tools requirements

## Config

Support simple configuration via a `.tlint.yml` file. This can be a global config in the home dir
or a per project config. Per project config will be prefered to global.

- goimports       - none
- go fmt          - none
- gometalinter    - .gometalinter.json in cwd
                    In the format of the config struct, which will be passed as the --config option
                    https://github.com/alecthomas/gometalinter/blob/master/config.go
- hadolint        - as hadolint does not support a config file config value may be placed in the .tlint.yml file

    ```
hadolint:
    ignore:
        - DL3007
```
- shellcheck      - use annotations in your scripts or include ignores in the config file
    ```
shellcheck:
    ignore:
        - SC2034
```
- checkmake       - .checkmake.ini in cwd
                    https://github.com/mrtazz/checkmake/blob/master/fixtures/exampleConfig.ini
                    https://github.com/mrtazz/checkmake/blob/master/man/man1/checkmake.1.md