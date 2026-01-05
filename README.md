# yatch

A tool for patching YAML files while preserving formatting and comments.

## Installation

```bash
go install github.com/alex-shpak/yatch/cmd/yatch@latest
```

## CLI Usage

```bash
# Update value from stdin to stdout
yatch '$.image.tag' 'v1.2.3' < input.yaml > output.yaml

# Update files in-place
yatch -in values.yaml -out values.yaml '$.version' '2.0'

# Update a comment instead of the value
yatch -in values.yaml --comment '$.service.port' 'Exposed port'

# Pipe from other commands
cat deployment.yaml | yatch '$.spec.replicas' '3'
```

### Options

- `-in` - Input YAML file (default: stdin)
- `-out` - Output YAML file (default: stdout)
- `--comment` - Update comment attached to the node instead of the node value

## Library Usage

```bash
go get github.com/alex-shpak/yatch/lib/yatch
```

```go
package main

import (
    "os"
    "github.com/alex-shpak/yatch/lib/yatch"
)

func main() {
    file, err := yatch.NewFile(reader)
    if err != nil {
        panic(err)
    }

    // Patch a value
    err = file.Patch("$.spec.image.tag", "v1.2.3")
    if err != nil {
        panic(err)
    }

    // Patch a comment
    err = file.PatchComment("$.spec.image.tag", "{\"$imagepolicy\": \"flux-system:my-policy\"}")
    if err != nil {
        panic(err)
    }

    // Get updated content
    updated := file.Content()
}
```

## How It Works
`yatch` parses YAML using `goccy/go-yaml` library to an AST, locates the target node using path expression, and performs byte-level replacement to preserve the original formatting, indentation, comments.
