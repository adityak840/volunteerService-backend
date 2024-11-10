include .env

up:
	@echo "Starting containers..."
	sudo docker-compose up --build -d --remove-orphans

down:
	@echo "Stoping containers..."
	sudo docker-compose down

build:
	go build -o ${BINARY} ./api/

start:
	@env MONGO_DB_USERNAME=${MONGO_DB_USERNAME} MONGO_DB_PASSWORD=${MONGO_DB_PASSWORD} ./${BINARY} 

restart: build start 

format_all_code:
	go fmt ./...

# connection string
# mongodb://admin:password@localhost:27017/volunteerService-backend-db?authSource=admin&readPreference=primary&appname=MongDB%20Compass&directConnection=true&ssl=false
