# Installation

## Build from source

Install the go toolchain: <https://golang.org/doc/install>

```bash
go install github.com/jdstrand/language-checker@latest

# Or pin a specific version
go install github.com/jdstrand/language-checker@v0.9.0
```

## Docker

You can run `language-checker` within docker. You will need to mount a volume that contains your source code and/or rules.

```bash
## Run with all defaults, within the mounted /src directory
docker run -v $(pwd):/src -w /src jdstrand/language-checker

## Provide rules config
docker run -v $(pwd):/src -w /src jdstrand/language-checker \
  language-checker -c my-rules.yaml
```
