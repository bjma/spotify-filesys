build:
	go build -o spfs ./main.go

run:
	#go run ./main.go This is a bit more annoying
	make build && ./spfs