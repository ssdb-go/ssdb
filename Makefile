PACKAGE_DIRS := $(shell find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | sort)

test: testdeps
	go test ./...
	go test ./... -short -race
	go test ./... -run=NONE -bench=. -benchmem
	env GOOS=linux GOARCH=386 go test ./...
	go vet
	cd internal/customvet && go build .
	go vet -vettool ./internal/customvet/customvet

testdeps: testdata/ssdb/ssdb-master

bench: testdeps
	go test ./... -test.run=NONE -test.bench=. -test.benchmem

.PHONY: all test testdeps bench

testdata/ssdb:
	mkdir -p $@
	wget --no-check-certificate https://github.com/ideawu/ssdb/archive/master.zip -O $@/master.zip
	cd $@ && unzip master

testdata/ssdb/ssdb-master: testdata/ssdb
	cd $@ && make all && ./ssdb-server -d ssdb.conf -s start

fmt:
	gofmt -w -s ./
	goimports -w  -local github.com/ssdb-go/ssdb ./

go_mod_tidy:
	set -e; for dir in $(PACKAGE_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go get -u ./... && \
	    go mod tidy -compat=1.17); \
	done
