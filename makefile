.PHONY: generate run tools dbml sqlc openapi migrate db-reset

tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/ogen-go/ogen/cmd/ogen@latest
	npm install --save-dev @dbml/cli @redocly/cli

db/schema.sql: db/schema.dbml
	npx dbml2sql db/schema.dbml -o db/schema.sql --postgres

api/bundle.yaml: api/openapi.yaml api/paths/*.yaml api/components/*.yaml
	npx redocly bundle api/openapi.yaml -o api/bundle.yaml

generate: db/schema.sql api/bundle.yaml
	go generate ./...
	go mod tidy

migrate:
	docker exec -i portfolio-manager-db psql -U postgres -d portfolio_manager < db/schema.sql

db-reset:
	docker exec -i portfolio-manager-db psql -U postgres -d postgres -c "DROP DATABASE IF EXISTS portfolio_manager WITH (FORCE);"
	docker exec -i portfolio-manager-db psql -U postgres -d postgres -c "CREATE DATABASE portfolio_manager;"
	make migrate

run: generate
	go run cmd/app/main.go
