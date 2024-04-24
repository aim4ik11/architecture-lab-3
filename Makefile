default: out/example

clean:
	rm -rf out

test:
	go test -timeout 30s github.com/aim4ik11/architecture-lab-3/painter

out/example:
	mkdir -p out
	go build -o out/example ./cmd/painter/main.go