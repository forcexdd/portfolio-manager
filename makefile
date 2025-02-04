.PHONY: open_db close_db run_web close_web

open_db :
	@docker compose -f src/deployments/docker-compose.yml up -d

close_db :
	@docker compose -f src/deployments/docker-compose.yml down

run_web : open_db
	@go run src/cmd/app/web.go

close_web: close_db
