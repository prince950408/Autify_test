
# Backend Engineer Take Home Test

A command line program that can fetch web pages and saves them to disk for later retrieval and browsing.

## Deployment

To install this project

``` bash
    go mod init fetch
    go mod tidy
    go build fetch
```

To test this project

```bash
    ./fetch https://www.google.com https://autify.com
    ./fetch --metadata https://www.google.com
```


When we use docker

```bash
    docker build -t fetch .
    docker run fetch https://www.google.com https://autify.com
    docker run fetch --metadata https://www.google.com
```

