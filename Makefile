APP=midnight
MODULE := github.com/branislavlazic/midnight
VERSION := v0.1

.PHONY: clean bin

all: clean swag frontend zip

swag:
	swag init -g main.go

clean:
	rm -rf bin release

frontend:
	cd webapp && yarn && yarn build

zip: release/$(APP)_$(VERSION)_linux_x86_64.tar.gz release/$(APP)_$(VERSION)_osx_x86_64.tar.gz

binaries: binaries/linux_x86_64/$(APP) binaries/osx_x86_64/$(APP)

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
