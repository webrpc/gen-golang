generate:
	webrpc-gen -schema=example.ridl -target=../../ -pkg=main -server -client -out=./example.gen.go
	goimports -w ./example.gen.go
