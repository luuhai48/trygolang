## Generate swagger document

```bash
go run ./src swag
```

## Create migration file

```bash
go run ./src migration make migrate_file_name
```

## Run migrations

```bash
go run ./src migration migrate
```

### Start web server

```bash
go run ./src runserver
# or
go run ./src run
```
