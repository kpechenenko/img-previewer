build:
	go build -o bin/image-previewer cmd/image-previewer/*

run:
	export SERVER_ADDR=":8080"; \
	build \
	./bin/image-previewer

test:
	go test ./... -v


#http://127.0.0.1:8080/preview/100/300/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg