# Yatch (yaml patch)
Yatch is a cli tool and a library to patch yaml files without editing its formatting

# Usage
## Usage as CLI
```
yatch [-in] [-out] '.path.to' 'qwerty'
```

## Usage as library
```shell
go get github.com/alex-shpak/yatch
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

