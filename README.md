# Jadwalin Auth Events Integrator

## Requirements
1. Go 1.18
2. Mongo
3. Redis

## Using the application
1. docker compose up -d 
2. go run cmd/main.go
3. docker compose down (if u need to stop the container)

##### If you have make command in your computer
1. make up
2. make run
3. make down

## default config.json
```
{
  "server": {
    "port": 8080
  },
  "database": {
    "connection_string": "mongodb://localhost:27017",
    "username": "admin",
    "password": "bantengmerah"
  },
  "cache": {
    "host": "localhost:6379",
    "password": "bantengmerah",
    "db": 0
  }
}
```
