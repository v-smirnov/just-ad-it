## Test app
App will send requests to the given list of URLs.  

### Implementation notes

Main logic can be found in ```internal/app/service/service.go```. ```Limiter``` chanel is used for limiting amount of parallel
requests, waitGroup allows waiting for all goroutines to be finished, ```responses``` channel is used for storing results
and fetching them later.
File```internal/app/infrastructure/client.go``` is used as simple http client implementation. Main method there returns
response body as a slice of bytes.

### How to run

First option is to run the app using docker:
- docker should be installed on your machine
- to build an image execute ```docker build . -t just-ad-app```
- to run an image execute ```docker run just-ad-app go run cmd/main.go -parallel 5 adjust.com google.com facebook.com test.wrong yandex.com twitter.com```

Second option is to run app locally:
- *Golang* should be installed on you machine
- from the application root folder execute ```go run -race cmd/main.go -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com```

After running the app you should be able to see some output information

### Running tests

To run the tests you can use the following command:

- from project root ```go test -race -v ./internal/app/*```