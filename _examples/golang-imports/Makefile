generate:
	webrpc-gen -schema=./proto/api.ridl -target=../../ -out=./api.gen.go -Pkg=main -Server -Client
	goimports -w ./api.gen.go
