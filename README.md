# PasetMe Backend

Using `Gin` and `Gorm`.

## API

| Method | URL | Status Code | Return | Description |
| --- | --- | --- | --- | --- |
| GET | /:token | 200 | content | token = key[,password] |
| POST | / | 201 | {status: 201, key: } | Create a permanent paste |
| POST | /once | 201 | {status: 201, key: } | Create a read once paste |
| PUT | /:key | 201 | {status: 201, key: } | Create a temporary paste |
| DELETE | /:key | 200 | {status: 200} | Delete a permanent key (not ready yet) |

### Creation JSON Format

```json
{
  "lang": "",
  "content": "",
  "password": ""
}
```

### Response JSON Format

#### Normal Response

```json
{

}
```

### Example

```bash
$ curl api.pasteme.cn/100 # access 100

$ curl api.pasteme.cn/101,123456 # access 101 with password 123456

$ curl -H 'Content-Type: application/json' \
-d '{"lang":"bash","content":"echo Hello"}' \
api.pasteme.cn # create an permanent paste

$ curl -H 'Content-Type: application/json' \
-d '{"lang":"bash","content":"echo Hello"}' \
api.pasteme.cn/once # create an temporary paste with random key

$ curl -X PUT -H 'Content-Type: application/json' \
-d '{"lang":"bash","content":"echo Hello"}' \
api.pasteme.cn/hello # create an temporary paste with specific key
```

## Build

```bash
$ bash dep.sh
$ go build -o pastemed
```

## Test

This script will test all packages if there is no param.

```bash
bash gotest.sh [package name]
```
