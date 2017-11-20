// Defining the base handler for handlers processing the requests for the test services

package handler

import (
    "errors"
    
    "app"
)

type ITestServiceContext interface {   // an interface for the contexts defined in the app/context.go for the test service
    OK(resp []byte) error
    NotFound() error
}

type IHandler interface {  // the base handler interface for the handlers processing requests for the test services
    Process() error  // starting the processes in the handler
}

func NewHandler(ctx ITestServiceContext) (ih IHandler, err error) {  // return a handler instance according to the type of the context
    switch c := ctx.(type) {
    case *app.LocalServiceTestServiceContext:
        return NewLocalServiceHandler(c)
    case *app.ServiceChainTestServiceContext:
        return NewServiceChainHandler(c)
    default:
        err = errors.New("Unknown type of the inputing test service context!")
    }
    return nil, err
}