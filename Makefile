SRC = main.go codegen.go util.go parse.go

rib:
	go build $(SRC)
