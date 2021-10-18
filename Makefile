build:
	docker-compose build balance-service 

run:
	docker-compose up balance-service

migrate:
	migrate -path ./schema -database 'postgres://postgres:passwd@0.0.0.0:5432/postgres?sslmode=disable' up
