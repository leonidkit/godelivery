## Create
```
curl -X POST -d '{"request_id":1234,"format_type":"SF","format":"<hello>world</hello>"}' http://localhost:9000/create
```

## Read
```
curl -X GET http://localhost:9000/read/SF/1234
```