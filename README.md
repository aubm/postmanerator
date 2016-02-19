[![Build Status](https://travis-ci.org/aubm/postmanerator.svg?branch=master)](https://travis-ci.org/aubm/postmanerator)

## What is it ?

This a program that parses exported Postman collections to generate HTTP API documentations.

## Usage

`postmanerator -theme='bootstrap_default' -output='./doc.html' /path/to/your/collection.postman.json`

## Features

- Custom themes

## Installation

Just download the latest appropriate [release on Github](https://github.com/aubm/postmanerator/releases).

Or install it from sources ...

- Follow the instructions to [get golang installed](https://golang.org/doc/install) on your machine
- Get the package with go get `go get github.com/aubm/postmanerator`
- Install it `go install`

## Todos

- Write tests
- Create a more flexible system for themes
- Make the README a bit more sexy

## License

[MIT License](LICENSE.md)