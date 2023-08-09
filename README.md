# :arrow_up: GoUp

Update all direct and only required indirect dependencies in multiple modules (even recursively)

## :gear: Install

```shell
go install github.com/mymmrac/goup@latest
```

## :kite: Usage

Update all direct modules for project

```shell
goup
```

Update in multiple directories

```shell
goup dir1 dir2 dir3
```

Update modules recursively in all directories and subdirectories

```shell
goup -r dir
```

Update in all subdirectories recursively excluding some that match pattern

```shell
goup -r -e "some*" -e "?other" dir
```

> For pattern matching, see [more here](https://pkg.go.dev/path/filepath#Match)

## :closed_lock_with_key: License

GoUp is distributed under [MIT licence](LICENSE)
