// Package strictenums exercises the -strictEnums option: an unknown enum name
// decodes to a reserved invalid value instead of the enum's zero value.
package strictenums

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=strictenums.ridl -target=../../../../gen-golang -pkg=strictenums -strictEnums -out=./strictenums.gen.go
