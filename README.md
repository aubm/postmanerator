[![Build Status](https://travis-ci.org/aubm/postmanerator.svg?branch=master)](https://travis-ci.org/aubm/postmanerator)

## What is it ?

This a program that parses exported Postman collections to generate HTTP API documentations.

## Usage

For now, you need [to have Go installed](https://golang.org/doc/install) to use it.

- Get the package with go get `go get github.com/aubm/postmanerator`
- Install it `go install`
- Use it `postmanerator -theme='bootstrap_default' -output='./doc.html' /path/to/your/collection.postman.json`

## Features

- Custom themes

## Todos

- Write tests
- Create a more flexible system for themes
- Make the README a bit more sexy

## License

[MIT License](LICENSE.md)