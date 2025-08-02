// gateway.go
package main

import (
	"github.com/nhdms/base-go/cmd/apis/api-gateway/internal"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/config"
	"github.com/nhdms/base-go/pkg/dbtool"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"log"
	"net/http"
)

func main() {
	service := app.NewAPI(nil)

	redis, err := dbtool.CreateRedisConnection(nil)
	if err != nil {
		log.Fatal("Failed to create redis connection: ", err)
	}

	proxy := internal.NewReverseProxy(redis, service.Options().Registry, nil)
	maxIdleConns := config.ViperGetIntWithDefault("http.max_idle_conns", 120)
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = maxIdleConns

	// Set up routes with authentication middleware
	httpHandler := http.NewServeMux()
	httpHandler.Handle("/", proxy)

	// Register handler
	service.Handle("/", httpHandler)
	service.Handle("/health-check", new(transhttp.HealthCheckHandler))

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
