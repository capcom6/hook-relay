// Package main Hook Relay webhooks delivery microservice
//
//	@title			Hook Relay API
//	@version		1.0.0
//	@description	Hook Relay is a webhook delivery microservice
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	Aleksandr Soloshenko
//	@contact.url	https://github.com/capcom6/hook-relay
//	@contact.email	i@capcom.me
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host			localhost:3000
//	@BasePath		/api/v1
package main

import (
	"runtime"
	"strconv"

	"github.com/capcom6/hook-relay/internal"
	"github.com/go-core-fx/healthfx"
	"github.com/samber/lo"
)

//go:generate swag init --parseDependency --outputTypes go -g ./main.go -o ./internal/server/docs

//nolint:gochecknoglobals // build metadata
var (
	appVersion   = "dev"
	appReleaseID = "0"
	appBuildDate = "unknown"
	appGitCommit = "unknown"
	appGoVersion = runtime.Version()
)

func main() {
	internal.Run(healthfx.Version{
		Version:   appVersion,
		ReleaseID: lo.Must1(strconv.Atoi(appReleaseID)),
		BuildDate: appBuildDate,
		GitCommit: appGitCommit,
		GoVersion: appGoVersion,
	})
}
