# App

## Build

```
$ GOOS=linux go build -o app .
```
```
$ docker build -t app:v2 .
```

## Run

```
$ docker run --rm \
  -e "APP_DATADIR=/var/lib/data" \
  -e "APP_HOST=203.0.113.10" \
  -e "APP_PORT=3306" \
  -e "APP_USERNAME=user" \
  -e "APP_PASSWORD=password" \
  -e "APP_DATABASE=test" \
  app:v2
```
