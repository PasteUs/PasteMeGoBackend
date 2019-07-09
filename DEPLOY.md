# Deployment

1. Download binary file `pastemed` in build branch
2. Execute the following command:

```bash
PASTEMED_DB_USERNAME= \ # Database username
PASTEMED_DB_PASSWORD= \ # Database password
PASTEMED_DB_SERVER= \ # MySQL server address
PASTEMED_DB_PORT= \ # MySQL server port
PASTEMED_DB_DATABASE= \ # Database name
./pastemed
```

## Params

| Name | Description | Example |
| :---: | --- | --- |
| PASTEMED_DB_USERNAME | Database username | `username` |
| PASTEMED_DB_PASSWORD | Database password | `password` |
| PASTEMED_DB_SERVER | MySQL server address | `localhost` |
| PASTEMED_DB_PORT | MySQL server port | `3306` |
| PASTEMED_DB_DATABASE | Database name | `pasteme` |
