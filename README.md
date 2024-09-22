# Browser Capabilities GoLang Project

This is a Go version of a library that matches user agents against the [Browscap](https://browscap.org/) database. It
uses a radix tree and backtracking to find the most specific pattern for a given user agent. The search index is
completely stored in memory. Once initialized, the full Browscap database requires about 20MB of RAM; however, during
initialization, it may require up to 170MB of RAM.

For browser data, it requires an external storage, which is currently implemented using SQLite, MySQL, or PostgreSQL.

## Installation

```bash
go get github.com/eugeniypetrov/browscap-go
```

## Usage

To build the cache, you need to download the PHP version of the database from [Browscap](https://browscap.org/) and run
the following code:


```go
db, err := sqlx.Open("sqlite3", "browscap.sqlite")
if err != nil {
    panic(err)
}

storage := browscap.NewSqliteBrowserStorage(db)
err := browscap.NewLoader(storage).Compile("full_php_browscap.ini")
if err != nil {
    panic(err)
}
```

It will take some time to parse the file and build the cache. Alternatively, you can install the CLI tool using the
following command:

```bash
go install github.com/eugeniypetrov/browscap-go@latest
```

and then run the following command:

```bash
browscap-go compile \
  -srorage=sqlite \
  -dsn=browscap.sqlite \
  -filename=full_php_browscap.ini
```

This will create an SQLite database with all the data from the full_php_browscap.ini file. After that, you can use the
following code to match user agents:

```go
bc, err := browscap.NewLoader(storage).Load()
if err != nil {
    panic(err)
}

browser, _ := bc.GetBrowser(userAgent)
```