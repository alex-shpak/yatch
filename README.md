# Yatch (YAML patch)
A CLI and a library to patch YAML files while preserving formatting.

# Usage
## Usage as CLI
```
go install github.com/alex-shpak/yatch/cmd/yatch@latest
yatch [-in=file] [-out=file] '.spec.image.tag' 'qwerty'
```

## Usage as library
```shell
go get github.com/alex-shpak/yatch/lib/yatch
```

```go
/* Read file or stdin above */
file, err := yatch.NewFile(reader)
if err != nil {
  return err
}

err := file.Patch(".spec.image.tag", "0.0.1")
if err != nil {
  return err
}

yamlfile.Content()
```

