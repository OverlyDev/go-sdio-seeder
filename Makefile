APP_PATH = ./cmd/sdio-seeder
APP_NAME = sdio-seeder

DOCKER_NAME = overlydev/sdio-seeder
DOCKER_TAG = latest
DOCKER_IMAGE = $(DOCKER_NAME):$(DOCKER_TAG)

clean:
	rm -rf build data downloads docker/sdio-seeder*

format:
	go fmt ./...

build: format
	CGO_ENABLED=0 go build -ldflags "-s -w" -o build/$(APP_NAME)-linux-x64 $(APP_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/$(APP_NAME)-win-x64.exe $(APP_PATH)
	ls -lh ./build

run: format
	go run $(APP_PATH)

image: build
	cp build/$(APP_NAME)-linux-x64 docker/sdio-seeder
	cd docker; docker build -t $(DOCKER_IMAGE) .

run_image: image
	docker run --rm -it $(DOCKER_IMAGE)