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
migrate -url postgres://user@host:port/database -path ./db/migrations -timeout 10 up
migrate -url postgres://user@host:port/database -path ./db/migrations -timeout 10 up 1
migrate -url postgres://user@host:port/database -path ./db/migrations -timeout 10 down
migrate -url postgres://user@host:port/database -path ./db/migrations -timeout 10 down 1
migrate help # for more info
```

## How to contribute

1. Fork the repo on Github.
2. Clone the `wallester/migrate` repo. Next steps are to be done in the `wallester/migrate repo` (not the fork).
3. Add your fork as a new remote (`git remote add myfork url-of-myfork`).
4. Create a new branch, do your work and commit the changes as usual.
5. Push your new branch to your fork (`git push myfork mybranch`).
6. Open a pull request in Github as usual.
