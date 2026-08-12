webrpc-gen Golang templates
===============================

This repo contains the templates used by the `webrpc-gen` cli to code-generate
webrpc Go server and client code.


## Usage

```
webrpc-gen -schema=example.ridl -target=golang -out=./example.gen.go -pkg=main -server -client

# or 
webrpc-gen -schema=example.ridl -target=github.com/webrpc/gen-golang@v0.6.0 -out=./example.gen.go -pkg=main -server -client

# or
webrpc-gen -schema=example.ridl -target=./local-go-templates-on-disk -out=./example.gen.go -pkg=main -server -client
```

As you can see, the `-target` supports default `golang`, any git URI, or a local folder :)

### Set custom template variables
Change any of the following values by passing `-option="Value"` CLI flag to `webrpc-gen`.

| webrpc-gen -option    | Default   | Description                                                                 | Added in |
|-----------------------|-----------|-----------------------------------------------------------------------------|----------|
| `-pkg=<name>`         | `"proto"` | package name                                                                | v0.5.0   |
| `-client`             | `false`   | generate client code                                                        | v0.5.0   |
| `-server`             | `false`   | generate server code                                                        | v0.5.0   |
| `-types=false`        | `true`    | don't generate types                                                        | v0.13.0  |
| `-json=sonic`         |           | use [sonic](https://github.com/bytedance/sonic) for JSON encoding           | v0.18.0  |
| `-json=jsoniter`      |           | use [jsoniter](https://github.com/json-iterator/go) for JSON encoding       | v0.12.0  |
| `-json=<pkg>`         |           | use alternative drop-in replacement import path for JSON encoding package   | v0.18.0  |
| `-fixEmptyArrays`     | `false`   | serialize lists per the schema, never as `null` (see Go [#27589][go27589])  | v0.13.0  |
| `-errorStackTrace`    | `false`   | enables error stack traces                                                  | v0.14.0  |
| `-webrpcHeader=false` | `true`    | enable client send webrpc version in http headers                           | v0.16.0  |
| `-schemaHash=false`   | `true`    | don't emit schema hash + version helper funcs (avoids merge conflicts)      | v0.30.0  |

Example:
```
webrpc-gen -schema=./proto.json -target=golang -out server.gen.go -pkg=main -server
```

## Empty arrays

A nil Go slice or map serializes as `null`, which does not match a schema that
says the field is a list or a map. `-fixEmptyArrays` makes the generated server
serialize collections the way the schema declares them:

| Schema field                | Go value             | JSON      |
|-----------------------------|----------------------|-----------|
| `- tags: []string`          | `nil`                | `[]`      |
| `- tags: []string`          | `[]string{}`         | `[]`      |
| `- tags?: []string`         | `nil`                | *absent*  |
| `- tags?: []string`         | `[]string{}`         | `[]`      |
| `- counts: map<string,int>` | `nil`                | `{}`      |
| `- counts?: map<string,int>`| `nil`                | *absent*  |

Required collections are always a collection. Optional ones keep all three
states, so a client can tell "the server said nothing" apart from "the server
said empty".

Nesting is walked all the way down, so the lists inside `[][]string`,
`[]Item`, `map<string,[]string>` and `map<string,Item>` get the same treatment.

The flag generates an `initNilSlices()` method on each schema struct and tags
optional fields `omitzero`. A few notes:

- `omitzero` needs Go 1.24+ in the module that consumes the generated code.
  Older toolchains ignore the tag and keep emitting `null`, as they do today.
- A nil **struct** is left as `null`. A zero struct is a different value rather
  than an empty collection, so filling one in would invent data and hide the bug
  that produced the nil.
- Fields that pin their own type with `go.field.type` are left alone unless that
  type is a slice, since the generator cannot know what an empty value means for
  an arbitrary type. This is what keeps an empty `json.RawMessage`, which is not
  valid JSON, from breaking the response.
- Fields that pin their own `go.tag.json` keep exactly the tag they asked for.
- A `[]byte` is base64-encoded by `encoding/json`, so it serializes as `""`
  rather than `[]`.
- The flag needs generated types, so it cannot be combined with `-types=false`
  or `-importTypesFrom`.

Generate the client with the same flag as the server, so both sides agree on the
struct tags.

## Set custom Go field meta tags in your RIDL file

| CLI option flag                              | Description                                                      |
|----------------------------------------------|------------------------------------------------------------------|
| `+ go.field.name = ID`                       | Set custom field name                                            |
| `+ go.field.type = uuid.UUID`                | Set custom field type (must be able to JSON unmarshal the value) |
| `+ go.type.import = github.com/google/uuid`  | Set custom field type's import path                              |
| `+ go.tag.json = id`                         | Set `json:"id"` struct tag                                       |
| `+ go.tag.db = id`                           | Set `db:"id"` struct tag                                         |

Example:
```ridl
struct User
  - ID: int64
    + go.tag.db = id
    + go.tag.json = id
  - UUID: string
    + go.field.type = uuid.UUID
    + go.type.import = github.com/google/uuid
    + go.tag.json = uuid
    + go.tag.db = uuid
  - Age: int
    + go.tag.db: age
  - Name: string
    + go.tag.db = name
  - PasswordHash: string
    + go.tag.db = passwd_hash
```

will result in

```go
import "github.com/google/uuid"

type User struct {
	ID           int64     `json:"id" db:"id"`
	UUID         uuid.UUID `json:"uuid" db:"uuid"`
	Name         string    `db:"name"`
	Age          int       `db:"age"`
	PasswordHash string    `json:"-" db:"passwd_hash"`
}
```

## Examples

See [_examples](./_examples)

## LICENSE

[MIT LICENSE](./LICENSE)

[go27589]: https://github.com/golang/go/issues/27589