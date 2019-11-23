# Deployment

## Quick start

> Version 3

1. Download binary file `pastemed` in build branch
2. Edit the `config.json`, `config.example.json` is the sample config file
3. Execute the `pastemed`

```bash
./pastemed -c ./config.json -d ./
```

## Param

```bash
Usage of pastemed:
  -c string
        -c <config file> (default "./config.json")
  -d string
        -d <data dir> (default "./")
  -debug
        --debug Using debug mode
  -version
        --version Print version information
```

## config.json

| Name | Description | Example |
| :---: | --- | --- |
| `address` | Listen address | `0.0.0.0` |
| `port` | Listen port | `8000` |
| `database.type` | Data Source, `mysql` for using MySQL, another for using SQLite3 | `mysql` |
| `database.username` | Database username | `username` |
| `database.password` | Database password | `password` |
| `database.server` | MySQL server address | `localhost` |
| `database.port` | MySQL server port | `3306` |
| `database.database` | Database name | `pasteme` |

## Upgrade from 2.x

See [Migrate Documentation](./MIGRATE.md)
