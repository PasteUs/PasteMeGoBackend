# Database Migrate

## Migrate from 2.x to 3.x

1. Download [db_transfer](https://github.com/PasteUs/PasteMeGoBackend/raw/build/db_transfer)
2. Execute the following command:

```bash
PASTEMED_DB_USERNAME= \
PASTEMED_DB_PASSWORD= \
PASTEMED_DB_SERVER= \
PASTEMED_DB_PORT= \
PASTEMED_DB_DATABASE= \
./db_transfer
```

### Params

| Name | Description | Example |
| :---: | --- | --- |
| PASTEMED_DB_USERNAME | Database username | `username` |
| PASTEMED_DB_PASSWORD | Database password | `password` |
| PASTEMED_DB_SERVER | MySQL server address | `localhost` |
| PASTEMED_DB_PORT | MySQL server port | `3306` |
| PASTEMED_DB_DATABASE | Database name | `pasteme` |

### Hint

When migration happend, there will be two new table named `permanents` and `temporaries` in your database, they are for 3.0 backend

Another thing you should know is that old data still valid for 2.x backend, migration function would not remove them, you can do it manual, or just let it go

## Migrate from 1.x to 2.x

Just do nothing

## Migrate from 0.9 to 1.1

Put [trans_0.9_to_1.1.php](https://github.com/LucienShui/PasteMe/blob/dbTrans/trans_0.9_to_1.1.php) into directory `web_root/lib`, then execute `php trans_0.9_to_1.1.php`

After upgrade, upgrade `config.php` by `config.example.php` is required
