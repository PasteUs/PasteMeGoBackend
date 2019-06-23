# PasetMe Backend

Using `Gin` and `Gorm`.

## API

| Method | URL | Status Code | Return | Description |
| --- | --- | --- | --- | --- |
| GET | /:token | 200 | content | token = key[,password] |
| POST | /*key | 201 | {status: 201, key: } | Add an permanent paste |
| DELETE | /:key | 200 | {status: 200} | Not yet |

## Build

```bash
$ bash dep.sh
$ go build -o pasteme_backend
```

## Test

This script will test all packages if there is no param.

```bash
./gotest [package name]
```
