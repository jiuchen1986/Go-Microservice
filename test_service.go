package main

import (
	"time"
    
    "app"
    "handler"
//    "github.com/jiuchen1986/Go-Microservice/app"
//    "github.com/jiuchen1986/Go-Microservice/handler"
	"github.com/goadesign/goa"
)

// TestServiceController implements the TestService resource.
type TestServiceController struct {
	*goa.Controller
    Delay time.Duration
}

// NewTestServiceController creates a TestService controller.
func NewTestServiceController(service *goa.Service, delay int64) *TestServiceController {
	return &TestServiceController{Controller: service.NewController("TestServiceController"), Delay: time.Duration(delay)}
}

// LocalService runs the local_service action.
func (c *TestServiceController) LocalService(ctx *app.LocalServiceTestServiceContext) error {
	// TestServiceController_LocalService: start_implement

	// Put your logic here
    
    if h, err := handler.NewHandler(ctx); err != nil {
        return err
    } else {
        if e := h.Process(c.Delay); e != nil {
            return e
        }
    }

	// TestServiceController_LocalService: end_implement
	return nil
}

// ServiceChain runs the service_chain action.
func (c *TestServiceController) ServiceChain(ctx *app.ServiceChainTestServiceContext) error {
	// TestServiceController_ServiceChain: start_implement

	// Put your logic here
    
    if h, err := handler.NewHandler(ctx); err != nil {
        return err
    } else {
        if e := h.Process(c.Delay); e != nil {
            return e
        }
    }

	// TestServiceController_ServiceChain: end_implement
	return nil
}
