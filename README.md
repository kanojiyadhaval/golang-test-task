

```bash
git clone https://github.com/kanojiyadhaval/golang-test-task
cd golang-test-task
```

```bash
docker-compose up --build
```



```bash
curl -X POST http://localhost:8080/message \
-H "Content-Type: application/json" \
-d '{"sender": "Alice", "receiver": "Bob", "message": "Hello Bob!"}'
```


```bash
curl "http://localhost:8081/message/list?sender=Alice&receiver=Bob"
```





If you encounter a Redis connection issue, ensure the following:

```bash
docker run -it --rm \
  --network <network_name> \
  -v $(pwd):/usr/src/app \
  -w /usr/src/app \
  golang:1.20 \
  go run cmd/processor/main.go
```

