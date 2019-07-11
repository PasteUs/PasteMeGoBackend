# PasetMe Backend

Using `Gin` and `Gorm`.

## API

[API Documentation](./API.md)

## Deploy

[Deploy guidance](./DEPLOY.md)

## pastemectl

[pastemectl](./pastemectl.sh) is backend's controllor, written by shell script, when using [auto deploy script](./installer.sh), `pastemectl` would be installed by default

Get more: [pastemectl Document](./PASTEMECTL_DOCUMENT.md)

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
