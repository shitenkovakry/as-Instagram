# Practices

### How to use?

```sh
GOOS=linux GOARCH=amd64 go build -o ./build/app/myapp ./server/main.go
docker-compose up -d
```

then wait for the docker infra to launch.

### Tests

In order to create user, do this:

```
curl -i -X POST \
  --data-binary '{"name": "ondrys","email": "rianby64@mail.ru"}' \
  http://localhost:8080/api/v1/createUser
```

In order to read photos, do this:

```
curl -i -X GET \
  --data-binary '{}' \
  http://localhost:8080//api/v1/photos
```

In order to add comment, do this:
```
curl -i -X POST \
  --data-binary '{"user_id": 1,"photo_id":1, "comment":"this is so ugly photo" }' \
  http://localhost:8080/api/v1/addComment
```
