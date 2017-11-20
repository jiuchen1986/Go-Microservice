// Defining the handler processing the requests for the local service

package handler

import (
    "fmt"
    "os"
    "errors"
    "strings"
    
    "app"
    "types"
    "utils"
)

type LocalServiceHandler struct {   // the handler processing the requests for the local service
    Ctx *app.LocalServiceTestServiceContext
}

func NewLocalServiceHandler(ctx *app.LocalServiceTestServiceContext) (h *LocalServiceHandler, err error) {  // generate a handler
    return &LocalServiceHandler{ctx}, nil
}

func (h *LocalServiceHandler) Process() error {  // the main requests process of the handler
    cha := make([]*types.ServiceStatus, 1)
    var er error
    cha[0], er = GetLocalServiceStatus() 
    if er != nil {
        h.Ctx.NotFound()
        return er
    }    
    resp := &types.TestServiceResponse{"1", cha}
    
    if strings.Compare(h.Ctx.SvcLo, resp.Chain[0].ServName) != 0 {
        return h.Ctx.NotFound()
    } else {
        resp_b, err := types.RespEncode(resp)
        if err != nil {
            h.Ctx.NotFound()
            return err
        } else {
            
            fmt.Printf("Send a response with OK!\n")
            fmt.Printf("Response body: \n")
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