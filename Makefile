APP=midnight
MODULE := github.com/branislavlazic/midnight
VERSION := v0.1

.PHONY: clean bin test

all: clean swag test zip

swag:
	swag init -g main.go

test: frontend
	docker-compose up -d postgres_test
	go test -count=1 -cover -v ./...
	docker-compose stop postgres_test

clean:
	rm -rf bin release

frontend:
	cd webapp && yarn && yarn build

zip: frontend release/$(APP)_$(VERSION)_linux_x86_64.tar.gz release/$(APP)_$(VERSION)_osx_x86_64.tar.gz

binaries: frontend binaries/linux_x86_64/$(APP) binaries/osx_x86_64/$(APP)

release/$(APP)_$(VERSION)_linux_x86_64.tar.gz: binaries/linux_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -C bin/linux_x86_64 $(APP)

binaries/linux_x86_64/$(APP):
	GOOS=linux GOARCH=amd64 go build -o bin/linux_x86_64/$(APP) main.go

release/$(APP)_$(VERSION)_osx_x86_64.tar.gz: binaries/osx_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -C bin/osx_x86_64 $(APP)

binaries/osx_x86_64/$(APP):
	GOOS=darwin GOARCH=amd64 go build -o bin/osx_x86_64/$(APP) main.go
