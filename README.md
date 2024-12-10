# Go Postgres (Wrong DDD)
The wrong dependency injection implementation in Go

## Features
- **Language**: [golang](https://golang.org/doc/install#install)
- **Framework**: [gin-gonic](https://github.com/gin-gonic/gin)
- **Database**: postgresql
- **Deploy**: docker, docker-compose
- **Other**: [ORM](https://github.com/jinzhu/gorm), [EnvoyProxy](https://www.envoyproxy.io/)

## [Structure](https://github.com/golang-standards/project-layout)
```
.
├── build
│    └── deploy                  # Deployment files. Docker, etc.
├── cmd                          
│    └── api                     # Starting point
├── deployments
│    ├── artifacts               # Postgres, etc.
│    └── etc                     # Environment files
├── pkg                          # App logic
└── README.md
```

## Environment variables
Name | Value
------------ | -------------
POSTGRES_USER|user
POSTGRES_PASSWORD|pass
POSTGRES_ADDRESS|postgres:5432
POSTGRES_DB|mydb
POSTGRES_URI|postgresql://user:pass@postgres:5432/mydb?sslmode=disable
CASBIN_MODEL_PATH|/artifacts/casbin/model.conf
CASBIN_POLICY_PATH|/artifacts/casbin/policy.csv

## Getting started
```shell script
# Build
docker build -t api -f build/deploy/Dockerfile .

# Deploy
cd deployments && docker-compose up -d
```
