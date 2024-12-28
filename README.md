# URL Shortener

An simple URL shortener built with golang and designed for ease of use and efficiency


## Tools and Technologies

| Tool/Technology Used| Purpose                                                                 |
|------------------|-------------------------------------------------------------------------|
| **Echo** | A high-performance, extensible, and minimal web framework for building Go web applications. |
| **SQLC**         | A tool for generating type-safe Go code from SQL queries, enabling efficient database interaction. |
| **PostgreSQL**   | A powerful, open-source relational database system for storing URL mappings and metadata. |
| **Redis**        | An in-memory data structure store used for caching and speeding up lookups. |
| **Migrate**      | A database migration tool for managing and versioning schema changes.  |
| **Viper**        | A configuration management tool in Go, used for handling application settings. |


## Build
```bash
# clone
$ git clone git@github.com:MashiroCl/UrlShortener.git
$ cd UrlShortener

# build
$ make up
$ go mod tidy
$ go build -o server main.go
$ ./server
```

## API Documentation
1. Create a Shortened URL
Endpoint: /api/url
Method: POST
Description: Accepts a request to shorten a URL, optionally allowing a custom code.

Request Body:
```json
{
  "original_url": "string",  // The original URL to be shortened (required)
  "custom_code": "string"    // Customized the URL shorten output(optional)
}
```
Response:
```json
{
  "short_url": "string",  // The generated shortened URL
  "expired_at": "time.Time" // shortened URL TTL
}
```


2. Redirect to the Original URL
Endpoint: /api/:code
Method: GET
Description: Redirects to the original URL based on the provided code.

**URL Parameters**:

| Parameter | Type   | Description                              |
|-----------|--------|------------------------------------------|
| `:code`   | String | The unique code corresponding to a URL. |

Example
Request:
```bash
GET /api/example123
```

Behavior:
- If the code exists: Redirects to the associated original URL (https://example.com).
- If the code does not exist: Returns an error response.
