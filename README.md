# :books: GoBookAPI :books: 

GoBookAPI is a sample RESTful JSON API written in Go, designed to manage books, provide user authentication/authorization
using tokens, and execute various operations including adding, deleting, viewing, and updating books. The API utilizes a PostgreSQL database for data storage and retrieval.

## Table of Contents

- [Features](#features)
- [File Structure](#structure)
- [Run Project Without Docker](#bare)
- [Run Project With Docker](#docker)
- [Usage Examples](#usage)
- [License](#license)


## Features <a name="features"></a>

- **System Monitoring:**
  - Logging
  - Metrics
  - Alerts
  - Healthcheck
  - Tracing
  - Error Handling

- **Security:**
  - CORS avoidance
  - IP-based rate limiting

- **Authentication/Authorization:**
  - Stateful token-based approach for secure authentication and authorization

- **Data Management:**
  - Book addition
  - Book deletion
  - Book updating
    - Partial data updates 
  - Book retrieval
    - Parsing and validating query string parameters
    - Listing data
    - Filtering and full-text search
    - Sorting and paginating
    - Returning pagination metadata
  - Optimistic concurrency control

- **Other:**
  - Continuous integration (CI)
  - Docker (Containerization)


## File Structure <a name="structure"></a>
```shell
GoBookAPI
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── api
│       ├── books.go
│       ├── context.go
│       ├── errors.go
│       ├── healthcheck.go
│       ├── helpers.go
│       ├── main.go
│       ├── middleware.go
│       ├── routes.go
│       ├── server.go
│       ├── tokens.go
│       └── users.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── data
│   │   ├── books.go
│   │   ├── filters.go
│   │   ├── models.go
│   │   ├── pages.go
│   │   ├── tokens.go
│   │   └── users.go
│   ├── jsonlog
│   │   └── jsonlog.go
│   └── validator
│       └── validator.go
└── migrations
    ├── 000001_create_books_table.down.sql
    ├── 000001_create_books_table.up.sql
    ├── 000002_add_books_check_constraints.down.sql
    ├── 000002_add_books_check_constraints.up.sql
    ├── 000003_add_books_indexes.down.sql
    ├── 000003_add_books_indexes.up.sql
    ├── 000004_create_users_table.down.sql
    ├── 000004_create_users_table.up.sql
    ├── 000005_create_tokens_table.down.sql
    └── 000005_create_tokens_table.up.sql
```


## Instructions
There are two options:
1) You can run the project without docker.
2) You can run the project with docker.

### 1. Run Project Without Docker<a name="bare"></a>

#### Requirements
- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [make](https://www.gnu.org/software/make/)
- [go-migrate](https://github.com/golang-migrate/migrate)
- [git](https://git-scm.com/downloads)
- Code Editor

#### Steps
- First, clone the repository:
```shell
git clone https://github.com/barantoraman/GoBookAPI.git
```
- Move into the directory:
```shell
cd GoBookAPI
```
- Create a PostgreSQL database with a named gobookapi.
- Create a POSTGRESQL superuser with named gobookapi.
- Set the password.
- Create the environment variables for databases:
```shell
export DB_DSN=YOUR_DATABASE_USER:YOUR_USER_PASS@/YOUR_DATABASE_NAME
```
- Edit the specified field in the Makefile that is related to the data source name, as per your requirement.

- Run the project:
```shell
make setup
```

- :tada: Now you can test the endpoints as demonstrated in the Usage Examples section. :tada:


### 2. Run Project With Docker<a name="docker"></a>

#### Requirements
- [Docker](https://docs.docker.com/get-docker/)
- [make](https://www.gnu.org/software/make/)
- [make](https://sourceforge.net/projects/gnuwin32/files/make/3.81/make-3.81.exe/download?use_mirror=nav&download=) (for windows)
- [git](https://git-scm.com/downloads)

#### Steps

- First, clone the repository:
```shell
git clone https://github.com/barantoraman/GoBookAPI.git
```
- Move into the directory:
```shell
cd GoBookAPI
```
- Build the containers  using the configuration in docker-compose.yaml:
```shell
make docker/build
```
- Start the containers using the configuration in docker-compose.yaml:
```shell
make docker/run
```
- Stop the containers:
```shell
make docker/stop
```


## Usage Examples<a name="usage"></a>

- Healthcheck:
  -  Endpoint to perform a health check and ensure the application is running smoothly.
```bash 
curl localhost:4000/v1/healthcheck
```


- Metrics:
  - Endpoint to retrieve application metrics for monitoring and analysis.
```bash 
curl localhost:4000/debug/vars 
```


- Signup:
  - If you wish to use the other functionalities of the application and test its additional endpoints, you need to sign up.
  Otherwise, you'll receive a (HTTP) 401 Unauthorized response.
```bash
curl -d '{
  "name": "Eliza Thornberry",
  "email": "ethornberry@example.com",
  "password": "random-password123"
}' localhost:4000/v1/users
```


- Login:
```bash
curl -d '{
  "email": "ethornberry@example.com",
  "password": "random-password123"
}' localhost:4000/v1/tokens/authentication
```


- Add a Book:
```bash 
curl -d '{
  "isbn": "0000000000001",
  "title": "Book1",
  "author": "Author1",
  "genres": ["genre1", "genre2", "genre3"],
  "pages": "100 pages",
  "language": "English",
  "publisher": "Publisher 1",
  "year": 1900
}' -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books"
```


- Update a Book:
  - Update all details of the book in the database.
```bash
curl -X PATCH -d '{
  "isbn": "0000000000001",
  "title": "Book1",
  "author": "Author1",
  "genres": ["genre1", "genre2", "genre3"],
  "pages": "100 pages",
  "language": "English",
  "publisher": "Publisher 1",
  "year": 1966
}' -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books/1"
```


- Partial Update: 
  -Updating specific fields of the book's information.
```bash
curl -X PATCH -d '{
  "year": 1902
}' -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books/1"
```


- Delete a Book by id:
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
-X DELETE "localhost:4000/v1/books/4"
```


- Get All Book:
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books"
```


- Get All Book:
  - For example, in this scenario, there will be a single results page with 2 results on that page.
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?cursor=1&cursor_size=2"
```


- Searching example: 
  - Searching for the desired book based on its features(Title + Author).
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?title=Book1&author=Author1"
```


- Searching example: 
  - Searching for the desired book based on its features(Title).
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?title=Book1" 
```


- Searching example:
  - Searching for the desired book based on its features(ISBN + Title + Author).
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?isbn=0000000000001&title=Book1&author=Author1"
```


- Searching example:
  - Searching for the desired book based on its features(ISBN).
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?isbn=0000000000001"      
```


- Searching example:
  - Searching for the desired book based on its features(Genres).
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?genres=genre1"
```


- Searching example:
  - Searching for the desired book based on its features(Title).
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?title=new+book"
```


- Full Text Search (enough to contain the word):
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?title=new"  
```


- Sort by ID (ASC):
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=id"
```


- Sort by ID (DESC): 
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=-id"
```


- Sort by Title (ASC):
```bash
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=title"
```


- Sort by Title (DESC):
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=-title"
```


- Sort by Year (ASC): 
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=year"
```


- Sort by Year (DESC):
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=-year"
```


- Sort by Author (ASC):
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=author"
```


- Sort by Author (DESC):
```bash 
curl -H "Authorization: Bearer UPXA2GVIOIE6KWUWDTT5NG4SH4" \
"localhost:4000/v1/books?sort=-author"
```