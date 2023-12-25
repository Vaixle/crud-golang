# crud-golang

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/inactive.svg)](https://www.repostatus.org/#inactive)
![GitHub License](https://img.shields.io/github/license/vaixle/crud-golang)


## Conntents:

- [Description](#Description)
- [Main Technologies](#Main-technologies)
- [Getting started](#Getting-started)
    - [Tests](#Tests)
    - [Start with Docker](#Start-with-docker)
    - [Swagger](#Swagger)



### Description
A simplified microservice on Golang to manage the task list (todo list).

---

### Main Technologies

| **Database** |       ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)         |
|:------------:|:--------------------------------------------------------------------------------------------------------------------------:|
| **Backend**  |                                                              ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)                                                              |
| **PaaS**  |        ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)        |

---

### Getting started
Run application.

#### Tests

```
go test -v -cover -race ./internal/...
```

#### Start with Docker

> Clone the repository

```
git clone https://github.com/Vaixle/empha-soft.git
```

> Run application in Docker container
```
docker-compose up
```

#### Swagger

```
http://localhost:8080/swagger/index.html#/
```
