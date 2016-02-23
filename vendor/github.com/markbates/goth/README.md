# Goth: Multi-Provider Authentication for Go [![Build Status](https://travis-ci.org/markbates/goth.svg)](https://travis-ci.org/markbates/goth)

Package goth provides a simple, clean, and idiomatic way to write authentication
packages for Go web applications.

Unlike other similar packages, Goth, lets you write OAuth, OAuth2, or any other
protocol providers, as long as they implement the `Provider` and `Session` interfaces.

This package was inspired by [https://github.com/intridea/omniauth](https://github.com/intridea/omniauth).

## Docs

The API docs can be found at [http://godoc.org/github.com/markbates/goth](http://godoc.org/github.com/markbates/goth)

## Installation

```text
$ go get github.com/markbates/goth
```

## Supported Providers

* Bitbucket
* Digital Ocean
* Dropbox
* Facebook
* GitHub
* Google+
* Lastfm
* Linkedin
* Spotify
* Twitch
* Twitter

## Examples

See the [examples](examples) folder for a working application that lets users authenticate
through Twitter, Facebook or Google Plus.

## Issues

Issues always stand a significantly better chance of getting fixed if the are accompanied by a
pull request.

## Contributing

Would I love to see more providers? Certainly! Would you love to contribute one? Hopefully, yes!

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Write Tests!
4. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request

## Contributors

* Mark Bates
* Tyler Bunnell
* Craig P Jolicoeur
* Zac Bergquist
* Ben Tranter
* Noah Shibley
* oov
* Tyler
* Rafael Quintela
* Noah
* DenSm
* HaiMing.Yin
* Masanobu YOSHIOKA
* bryanl
* dkhamsing
* Jerome Touffe-Blin
* Jonathan Hall
