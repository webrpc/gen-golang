webrpc-gen Golang templates
===============================

This repo contains the templates used by the `webrpc-gen` cli to code-generate
webrpc Go server and client code.


## Usage

```
webrpc-gen -schema=example.ridl -target=golang -pkg=main -server -client -out=./example.gen.go
```

or 

```
webrpc-gen -schema=example.ridl -target=github.com/webrpc/gen-golang@v0.6.0 -pkg=main -server -client -out=./example.gen.go
```

or

```
webrpc-gen -schema=example.ridl -target=./local-go-templates-on-disk -pkg=main -server -client -out=./example.gen.go
```

As you can see, the `-target` supports default `golang`, any git URI, or a local folder :)


## Examples

See [_examples](./_examples)

## LICENSE

[LICENSE](./LICENSE)
