# migrate
[![Travis](https://travis-ci.org/wallester/migrate.svg?branch=master)](https://travis-ci.org/wallester/migrate)
[![Go Report Card](https://goreportcard.com/badge/github.com/wallester/migrate)](https://goreportcard.com/report/github.com/wallester/migrate)

Command line tool for PostgreSQL migrations 

## Features

* Runs migrations in transactions
* Stores migration version details in auto-generated table ``schema_migrations``.

## Usage

```bash
migrate -url postgres://user@host:port/database -path ./db/migrations create add_field_to_table
migrate -url postgres://user@host:port/database -path ./db/migrations up
migrate -url postgres://user@host:port/database -path ./db/migrations up 1
migrate -url postgres://user@host:port/database -path ./db/migrations down
migrate -url postgres://user@host:port/database -path ./db/migrations down 1
migrate help # for more info
```
