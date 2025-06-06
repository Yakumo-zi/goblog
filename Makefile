.PHONY: build run test clean fmt vet lint deps dev docker-build docker-run docker-stop docker-logs docker-clean backup-test db-migrate db-check db-backup db-restore

build:
	@echo "Building..."
	go build -o bin/goblog cmd/server/main.go

run:
	@echo "Running server..."
	go run cmd/server/main.go

test:
	@echo "Running tests..."
	go test -v ./test/...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./test/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

clean:
	@echo "Cleaning..."
	rm -f bin/goblog
	rm -f coverage.out
	rm -f coverage.html

deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

ent-gen:
	@echo "Generating Ent code..."
	cd ent && go run generate.go

dev: clean ent-gen fmt vet build run

ci: fmt vet test

# Docker相关命令
docker-build:
	@echo "Building Docker image..."
	docker build -t goblog:latest .

docker-run:
	@echo "Running with Docker Compose..."
	docker-compose up -d

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f goblog

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker rmi goblog:latest 2>/dev/null || true

# 数据库相关
db-migrate:
	@echo "Running PostgreSQL migration..."
	./scripts/migrate_postgres.sh migrate

db-check:
	@echo "Checking PostgreSQL connection..."
	./scripts/migrate_postgres.sh check

db-backup:
	@echo "Creating database backup..."
	@mkdir -p backups
	docker exec goblog-postgres pg_dump -U goblog goblog > backups/goblog_backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup created in backups/ directory"

db-restore:
	@echo "Restoring database from backup..."
	@echo "Usage: make db-restore BACKUP_FILE=path/to/backup.sql"
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "Error: BACKUP_FILE parameter is required"; \
		echo "Example: make db-restore BACKUP_FILE=backups/goblog_backup_20240101_120000.sql"; \
		exit 1; \
	fi
	docker exec -i goblog-postgres psql -U goblog goblog < $(BACKUP_FILE)

# 备份相关
backup-test:
	@echo "Testing backup endpoint..."
	@echo "Make sure you have TOKEN environment variable set"
	@curl -H "Authorization: Bearer $$TOKEN" \
		-o "backup_$$(date +%Y%m%d_%H%M%S).zip" \
		"http://localhost:8080/api/articles/backup"

help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the server"
	@echo "  test           - Run unit tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  ent-gen        - Generate Ent code"
	@echo "  dev            - Development workflow"
	@echo "  ci             - CI workflow"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  docker-logs    - Show Docker logs"
	@echo "  docker-clean   - Clean Docker resources"
	@echo "  db-migrate     - Run PostgreSQL migration"
	@echo "  db-check       - Check PostgreSQL connection"
	@echo "  db-backup      - Create database backup"
	@echo "  db-restore     - Restore database from backup"
	@echo "  backup-test    - Test backup endpoint"
	@echo "  help           - Show this help" 