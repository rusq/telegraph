# Crippled Telegraph Library for Go

[![Build Status](https://travis-ci.com/rusq/telegraph.svg?branch=master)](https://travis-ci.com/rusq/telegraph)

![COME GET IT](./come_get_it.jpg)

This Go package provides some basic functions to interact with
[Telegraph](https://telegra.ph).

## Example

```go
package main

import (
    "log"
    "os"

    "github.com/rusq/telegraph"
)

func main() {
    f, err := os.Open("cat_pic.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    result, err := telegraph.Upload(context.Background(), f)
    if err != nil {
        log.Fatal("oh, man :(")
    }

    log.Printf("%v", result)
}
```

## Crippled?
An inquisitive reader might enquire:  "Why Crippled?"

Well, it's rude to leave questions unanswered, so I'll respond:
"Because currently, it only supports the methods listed below".

* Upload

## Contributing

You are more than welcome to add more methods or open an issue for me to add
one, which I cannot promise that I will do in a timely manner.

Most likely no one will ever find this repository anyway, so I don't have to
worry.
