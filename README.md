webrpc-gen Golang templates
===============================

This repo contains the templates used by the `webrpc-gen` cli to code-generate
webrpc Go server and client code.


## Usage

```
webrpc-gen -schema=example.ridl -target=golang -out=./example.gen.go -Pkg=main -Server -Client

# or 
webrpc-gen -schema=example.ridl -target=github.com/webrpc/gen-golang@v0.6.0 -out=./example.gen.go -Pkg=main -Server -Client

# or
webrpc-gen -schema=example.ridl -target=./local-go-templates-on-disk -out=./example.gen.go -Pkg=main -Server -Client
```

As you can see, the `-target` supports default `golang`, any git URI, or a local folder :)

### Set custom template variables
Change any of the following values by passing `-Option="Value"` CLI flag to webrpc-gen.

| CLI option flag      | Description                | Default value              |
|----------------------|----------------------------|----------------------------|
| `-Pkg=pkgname`       | package name               | `"proto"`                  |
| `-Client`            | generate client code       | unset (false)              | 
| `-Server`            | generate server code       | unset (false)              |

Example:
```
webrpc-gen -schema=./proto.json -target=golang -out openapi.gen.yaml -Pkg=main -Client -Server
```


## Examples

See [_examples](./_examples)

## LICENSE

[LICENSE](./LICENSE)
