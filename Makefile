# ==================================================================================== #
# HELPERS                                        
# ==================================================================================== #

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT 								   
# ==================================================================================== #

migrate:
	@echo "Applying migrations..."
	migrate -path=./migrations -database=${DB_DSN} up
	@echo "Migrations applied."
# DB_DSN='postgres://username:password@host/dbname'
# You can replace the placeholders (username, password, host, database name)
# with your actual database connection details and replace ${DB_DSN}
# - username: Your database username
# - password: Your database password
# - host: The host or address of your database server
# - dbname: The name of your database

# ==================================================================================== #
# PRODUCTION								           
# ==================================================================================== # 

#Run the project locally 
.PHONY: go/tidy 
go/tidy: 
	go mod tidy 

.PHONY: go/run
go/run:
	go run ./cmd/api

.PHONY: setup
setup: confirm go/tidy go/run migrate

# Build the containers  using the configuration in docker-compose.yaml
docker/build:
	docker compose build

# Start the containers using the configuration in docker-compose.yaml
docker/run: confirm
	docker compose up -d

# Stop the containers
docker/stop: confirm
	docker compose down

# ==================================================================================== #
# QUALITY CONTROL								       
# ==================================================================================== #
# audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy -compat=1.17
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor