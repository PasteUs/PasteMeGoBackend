# Deployment

> Version 3

1. Download binary file `pastemed` in build branch
2. Edit the `config.json`, `config.example.json` is the sample config file, if you don't know how, just execute `cp config.example.json config.json`
3. Execute the `pastemed`

```bash
./pastemed -c ./config.json
```

## config.json

> If you don't know what you are looking for, just make everything default

| Name | Description | Example |
| :---: | --- | --- |
| address | Listen address | `0.0.0.0` |
| port | Listen port | `8000` |
| debug | Whether use debug mod | `false` |
| database.type | Data Source, `mysql` for using MySQL, another for using SQLite3 | `mysql` |
| database.username | Database username | `username` |
| database.password | Database password | `password` |
| database.server | MySQL server address | `localhost` |
| database.port | MySQL server port | `3306` |
| database.database | Database name | `pasteme` |

## Upgrade from 2.x

See [Migrate Documentation](./MIGRATE.md)
