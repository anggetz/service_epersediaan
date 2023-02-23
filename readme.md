## QUICK START 

1. install go swagger bin from here 

```
go install github.com/swaggo/swag/cmd/swag@latest

```

2. ensure you already register path of golang bin folder in yours environment path, ex: ${GOROOT}/bin

try this commad in yours terminal before continue to next step

```
swag

```


3. setup yours .env.local 
4. run

```
go run main.go
```

#swagger

for generate swagger use this command below
```
swag init
```

note: every time you change the swaggo annotation please use swag init first, then you can run the service normally (go run main.go)

## about project structure

i am using DDD patterns for this project, 