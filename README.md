# server-common-go
This repository contains public packages to be used to build a golang server (`grpc` or `rest`).
it contains useful packages:
- `configuration` extracts configuration from `yaml` file using struct model annotation (uses [viper](https://github.com/spf13/viper))
- `log` wrapper for [logrus logger](https://github.com/sirupsen/logrus)
- `database` wrapper for [go-gorm/gorm package](https://github.com/go-gorm/gorm)
- `http` to build an http server, wrapper for [gin-gonic/gin package](https://github.com/gin-gonic/gin)
- multiple utils packages like `iso8601` duration or `crypto`

## Examples

Some `REST` and `grpc` server example are available in `api` directory.
