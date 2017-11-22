// Defining the handler processing the requests for the local service

package handler

import (
    "fmt"
    "os"
    "errors"
    "strings"
    "time"
    "strconv"
    "math/rand"

    "app"
    "types"
    "utils"    
//    "github.com/jiuchen1986/Go-Microservice/app"
//    "github.com/jiuchen1986/Go-Microservice/types"
//    "github.com/jiuchen1986/Go-Microservice/utils"
)

type LocalServiceHandler struct {   // the handler processing the requests for the local service
    Ctx *app.LocalServiceTestServiceContext
}

func NewLocalServiceHandler(ctx *app.LocalServiceTestServiceContext) (h *LocalServiceHandler, err error) {  // generate a handler
    return &LocalServiceHandler{ctx}, nil
}

func (h *LocalServiceHandler) Process(delay time.Duration) error {  // the main requests process of the handler
    fmt.Println("handler.localservhandler: Delay for ", delay)
    time.Sleep(delay)
    
    
    /*
    req_header := h.Ctx.RequestData.Request.Header
    fmt.Println("handler.localservhandler: Get headers from the request: ")
    for k, v := range req_header {
        fmt.Printf("%s: %v\n", k, v)
    }
    */
    
    sub_chains_resp := make([]*types.ServiceChain, 0)
    chain_resp := make([]*types.ServiceStatus, 1)
    var er error
    chain_resp[0], er = GetLocalServiceStatus() 
    if er != nil {
        h.Ctx.NotFound()
        return er
    }
    rand.Seed(315)    
    resp := &types.TestServiceResponse{&types.ServiceChain{chain_resp[0].ServName, 
                                                           strconv.Itoa(rand.Int()),
                                                           chain_resp,  
                                                           "1"}, sub_chains_resp}
    
    if strings.Compare(h.Ctx.SvcLo, resp.MainChain.Starter) != 0 {
        return h.Ctx.NotFound()
    } else {
        resp_b, err := types.RespEncode(resp)
        if err != nil {
            h.Ctx.NotFound()
            return err
        } else {
            
            fmt.Printf("handler.localservhandler: Send a response with OK!\n")
            fmt.Printf("handler.localservhandler: Response body: \n")
            fmt.Println(utils.Convert(resp_b))
            
            return h.Ctx.OK(resp_b)
        }
    }
    return nil    
}

func GetLocalServiceStatus() (st *types.ServiceStatus, err error) {
    st = &types.ServiceStatus{"1", "", ""}
    
    if e := os.Getenv("TEST_SERVICE_NAME"); strings.Compare(e, "") != 0 {
        st.ServName = e
    } else {
        return nil, errors.New("Env TEST_SERVICE_NAME is missing!")
    }
    if e := os.Getenv("TEST_SERVICE_VERSION"); strings.Compare(e, "") != 0 {
        st.Version = e
    } else {
        return nil, errors.New("Env TEST_SERVICE_VERSION is missing!")
    }
    return st, nil
}