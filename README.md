# Jadwalin Auth Events Integrator

## Requirements
1. Go 1.18
2. Mongo
3. Redis

## Using the application
1. docker compose up -d 
2. go mod tidy
3. go run main.go
4. docker compose down (if u need to stop the container)

##### If you have make command in your computer
1. make up
2. make run
3. make down

## Default config.json
```
{
  "server": {
    "port": 8080
  },
  "database": {
    "connection_string": "mongodb://localhost:27017",
    "name": "jadwalin",
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

## Debugging mode purpose in VSCODE (.vscode/launch.json)
```
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Application",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "args": [
        "run"
      ]
    }
  ]
}
```
