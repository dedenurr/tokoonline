# Online Shop Project
1. Jalankan docker untuk PostgresSQL

```
docker run --name postgresql -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=database -d -p 5433:5432 postgres:16
```

2. Download beberapa library
```
github.com/gin-gonic/gin
github.com/jackc/pgx/v5/stdlib
github.com/google/uuid
golang.org/x/crypto
```

3. Export Environment variable yang dibutuhkan
```
export DB_URI=postgres://user:password@localhost:5433/database?sslmode=disable
```

4. Jalankan Program
```
go run main.go
```