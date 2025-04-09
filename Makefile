
build:
	docker build -t language-checker -f dev/Dockerfile .

run:
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/jdstrand/language-checker \
		language-checker "./*.go ./**/*.go ./**/**/*.go *.yaml"

.PHONY: build run

prof:
	go test -bench=. -run=^$$ -cpuprofile cpu.prof -memprofile mem.prof ./cmd

prof-mem: prof
	pprof -top mem.prof | head -n 10

# pprof -http=localhost:8080 mem.prof

prof-cpu: prof
	pprof -http=localhost:8080 cpu.prof

test:
	go vet ./...
	go test -v ./...

.PHONY: prof prof-mem prof-cpu
