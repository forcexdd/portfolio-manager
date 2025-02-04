export GOPROXY = https://proxy.golang.org,direct
.PHONY: open_db close_db run_web close_web run_desktop close_desktop

open_db :
	@docker compose -f src/deployments/docker-compose.yml up -d

close_db :
	@docker compose -f src/deployments/docker-compose.yml down

run_web : open_db
	@go run src/cmd/app/*.go web

close_web: close_db

run_desktop : open_db
	go env
	@go run src/cmd/app/*.go desktop

close_desktop: close_db
