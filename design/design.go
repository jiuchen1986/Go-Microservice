// This a the definitions of the API for the testrest
// using the DSL providing by the goa, 
// see https://goa.design/design/overview/

package design

import (
    . "github.com/goadesign/goa/design"
    . "github.com/goadesign/goa/design/apidsl"
)

// As the goa API design language is a DSL implemented in Go and is not Go, 
// the “dot import” is used here. 
// The generated code or any of the actual Go code in goa 
// does not make use of “dot imports”. 

var _ = API("TEST REST API", func() {        // "TEST REST API" is the name of the API used in docs
    Title("TEST REST SERVICE")      // Documentation title
    Description("A restful service for test") // Longer documentation description
    Host("localhost:8081")                // Host used by Swagger and clients
    Scheme("http")                   // HTTP scheme used by Swagger and clients
    BasePath("/api")                  // Base path to all API endpoints
    Consumes("application/json")      // Media types supported by the API
    Produces("application/json")      // Media types generated by the API
    
})

var _ = Resource("TestService", func() {     // Defines the Test Service resource
    DefaultMedia("application/json")
    
    Action("local_service", func() {     // Defines the local service action
        Routing(GET("/:svcLo/"))         // The relative path to the local service endpoints
        Description("return the local service")
        Params(func() {
            Param("svcLo", String, "local service")  // Defines svcLo parameter as path segment
                                                    // captured by :svcLo
        })
        Response(OK)
        Response(NotFound)
    })

    Action("service_chain", func() {     // Defines the local service action
        Routing(GET("/:svcLo/*svcOther"))         // The relative path to the service chain endpoints
        Description("follow the service chain to the next service")
        Params(func() {
            Param("svcLo", String, "local service")
            Param("svcOther", String, "other services following the service chain")
        })
        Response(OK)
        Response(NotFound)
    })
})    

