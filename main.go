//go:generate goagen bootstrap -d design

package main

import (
    "app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("TEST REST API")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "Test Service" controller
	c := NewTestServiceController(service)
	app.MountTestServiceController(service, c)

	// Start service
	if err := service.ListenAndServe(":8082"); err != nil {
		service.LogError("startup", "err", err)
	}

}
