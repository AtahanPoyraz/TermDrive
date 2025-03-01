@echo off

REM Start Docker Compose setup
echo Starting Docker Compose setup...

REM Stop and remove all containers, and volumes defined in the docker-compose file
docker-compose down -v

REM Build and run services
docker-compose up --build -d

REM Check the status
echo Services are up and running:
docker-compose ps

REM Run go mod tidy to clean and update dependencies
echo Running go mod tidy...
go mod tidy

echo Go modules updated.

pause