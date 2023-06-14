# Practices

### How to use?

```sh
GOOS=linux GOARCH=amd64 go build -o ./build/app/myapp ./server/main.go
docker-compose up -d
```

then wait for the docker infra to launch.

### Tests
test address: `http://insta-test.mas3.co`

### Tests

In order to create user, do this:
```
curl -i -X POST \
  --data-binary '{"name": "ondrys","email": "rianby64@mail.ru"}' \
  http://localhost:8080/api/v1/users/register
```

In order to update email of user, do this:
```
curl -i -X PUT \
  --data-binary '{"user_id": 1, "new email":"kry@mail.ru"}' \
  http://localhost:8080/api/v1/users/update/email
```

In order to update name of user, do this:
```
curl -i -X PUT \
  --data-binary '{"user_id": 1, "new name":"kry"}' \
  http://localhost:8080/api/v1/users/update/name
```

In order to delete user, do this:
```
curl -i -X DELETE \
  --data-binary '{"user_id": 1}' \
  http://localhost:8080/api/v1/users/delete
```

In order to read user, do this:
```
curl -i -X GET \
  http://localhost:8080/api/v1/user/1
```



In order to read photos, do this:
```
curl -i -X POST \
  --data-binary '{"user_id": 1}' \
  http://localhost:8080/api/v1/photos
```

In order to read photo, do this:
```
curl -i -X GET \
  http://localhost:8080/api/v1/photo/1/1
```

In order to add photo, do this:
```
curl -i -X POST \
  --data-binary '{"user_id": 1}' \
  http://localhost:8080/api/v1/photos/add
```

In order to delete photo, do this:
```
curl -i -X DELETE \
  --data-binary '{"photo_id": 1}' \
  http://localhost:8080/api/v1/photos/delete
``




In order to add comment, do this:
```
curl -i -X POST \
  --data-binary '{"user_id": 1,"photo_id":1, "comment":"this is so ugly photo"}' \
  http://localhost:8080/api/v1/comments/add
```

In order to delete comment, do this:
```
curl -i -X DELETE \
  --data-binary '{"comment_id": 1}' \
  http://localhost:8080/api/v1/comments/delete
```

In order to update comment, do this:
```
curl -i -X PUT \
  --data-binary '{"comment_id": 1, "new comment":"kry is cute girl"}' \
  http://localhost:8080/api/v1/comments/update
```

In order to read comments, do this:
```
curl -i -X GET \
  http://localhost:8080/api/v1/comments/1/1
```




In order to add like, do this:
```
curl -i -X POST \
  --data-binary '{"user_id": 1, "photo_id":1}' \
  http://localhost:8080/api/v1/likes/add
```

In order to count likes, do this:
```
curl -i -X GET \
  http://localhost:8080/api/v1/likes/1
```

In order to delete like, do this:
```
curl -i -X DELETE \
  --data-binary '{"user_id": 1, "photo_id": 1}' \
  http://localhost:8080/api/v1/likes/delete
```

Do not forget to use MongoDB Compas to see the results.

curl 'http://localhost:8080/example-upload-file' \
-X 'POST' \
-H 'Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryBBfCeIeb5cVcTyaG' \
-H 'Content-Length: 37866' \
--data-binary $'------WebKitFormBoundaryBBfCeIeb5cVcTyaG\r\nContent-Disposition: form-data; name="file-upload"; filename=" \u0438\u0441\u043f21.jpg"\r\nContent-Type: image/jpeg\r\n\r\n\r\n------WebKitFormBoundaryBBfCeIeb5cVcTyaG--\r\n'
