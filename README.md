[![Build Status](https://travis-ci.org/aubm/postmanerator.svg?branch=master)](https://travis-ci.org/aubm/postmanerator)

## What is it ?

This a program that parses exported Postman collections to generate HTTP API documentations.

## Usage

`postmanerator -theme='/path/to/theme' -output='./doc.html' /path/to/your/collection.postman.json`

## Themes

The whole point of this tool is to be able to generate beautiful documentations from a Postman collection export.
This is done by exposing the collection data with a few helpers to a template.

The template is defined in a file named `index.tpl` and stored in a folder with the name of your choice. This
folder must be given to postmanerator using the `-theme` option (see the example above).

Some example templates are available [here](https://github.com/aubm/postmanerator/tree/master/ressources/themes).
Feel free to use one of these or create your own !

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