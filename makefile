.PHONY: run_web close_web

run_web :
	@docker compose -f src/deployments/docker-compose.yml up -d; \
	go run src/cmd/app/web.go

close_web:
	@docker compose -f src/deployments/docker-compose.yml down
