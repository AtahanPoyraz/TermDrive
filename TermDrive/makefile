include ./.env

.PHONY: all
.DEFAULT_GOAL := all

# 'all' target: Runs the setup and server start processes in sequence.
all: run-setup run-server

# 'run-setup' target: Runs the setup script based on the operating system.
# If the system is Windows, it runs setup.bat, otherwise it runs setup.sh.
run-setup:
ifeq ($(OS), Windows_NT)
	call ./setup/setup.bat
else
	bash ./setup/setup.sh
endif

# 'download-modules' target: Downloads and updates the Go modules.
download-modules:
	go mod tidy
	go mod download

# 'run-server' target: Runs the Go server by executing the main.go file.
run-server:
	cd server && go run main.go

# 'create-admim' target: Create admin user.
create-admin:
	cd server && go run main.go --option createadmin --username $(USERNAME) --email $(EMAIL) --password $(PASSWORD)

# 'connect-db' target: Connects to the PostgreSQL database inside the Docker container using psql command.
connect-db:
	docker compose exec -it ${CONTAINER_NAME} psql -U ${DATABASE_USER} -d ${DATABASE_NAME}
