//go:generate goagen bootstrap -d design

package main

import (
    "flags"
    
    "github.com/jiuchen1986/Go-Microservice/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

var delay int64  // The processing delay for this service

func init() {
    const (
        defaultDelay = 0
        usage = "the processing delay for this service"
    )
    flag.Int64Var(&delay, "delay", defaultDelay, usage)
    flag.Int64Var(&delay, "d", defaultDelay, usage + " (shorthand)")
}

func main() {
	flag.Parse()
    // Create service
	service := goa.New("TEST REST API")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "Test Service" controller
	c := NewTestServiceController(service, delay)
	app.MountTestServiceController(service, c)

	// Start service
	if err := service.ListenAndServe(":8082"); err != nil {
		service.LogError("startup", "err", err)
	}

}
