.PHONY: dev
dev:
	~/go/bin/air # & cd web && pnpm dev --host


.PHONY: build
build:
	go build -o main ./cmd/budgetto   
