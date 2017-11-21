// Defining the handler processing the requests for the service chain

package handler

import (
    "fmt"
    "strings"
    "net/http"
    "io/ioutil"
    "strconv"
    
    "app"
    "types"
    "utils"
)

var trace_header = [...]string{"X-Request-Id",
                               "X-B3-Traceid",
                               "X-B3-Spanid",
                               "X-B3-Parentspanid",
                               "X-B3-Sampled",
                               "X-B3-Flags",
                               "X-Bt-Span-Context"}  // Headers for distributed tracing
const SVC_TO_PORT = "8082"   // The port of the next service

type ServiceChainHandler struct {   // the handler processing the requests for the local service
    Ctx *app.ServiceChainTestServiceContext
}

func NewServiceChainHandler(ctx *app.ServiceChainTestServiceContext) (h *ServiceChainHandler, err error) {  // generate a handler
    return &ServiceChainHandler{ctx}, nil
}

func (h *ServiceChainHandler) Process() error {  // the main requests process of the handler
    cha, er := GetLocalServiceStatus()
    if er != nil {
        h.Ctx.NotFound()
        return er
    }        
    
    if strings.Compare(h.Ctx.SvcLo, cha.ServName) != 0 {
        return h.Ctx.NotFound()
    } else {
        resp_f, er := h.FollowChain()
        if er != nil {
            h.Ctx.NotFound()
            return er
        }
        l, er := strconv.Atoi(resp_f.Len)
        if er != nil {
            h.Ctx.NotFound()
            return er
        }
        resp_f.Len = strconv.Itoa(l + 1)
        cha.Order = resp_f.Len
        resp_f.Chain = append(resp_f.Chain, cha)        
        resp_b, err := types.RespEncode(resp_f)
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

func (h *ServiceChainHandler) FindNextService() (host, port string, err error) {  // Return the endpoint of the next service
    svcTo := strings.Split(h.Ctx.SvcOther, "/")[0]
    fmt.Println("The next service: ", svcTo)
    return svcTo, "8082", nil
    // return "10.0.2.15", "8082", nil
}

func PropTraceInfo(ih, oh *http.Header) error {  // Collect and progapate the headers from the incoming request to the outgoing request for tracing
    for _, h := range trace_header {
        if v, ok := (*ih)[h]; ok {
            fmt.Println("Found an incomming header: ", h, "=", v[0])
            if v[0] != "" {
                (*oh).Add(h, v[0])
                fmt.Println("Add a header: ", h, "=", v[0])
            }
        }
    }
    return nil
}

func (h *ServiceChainHandler) FollowChain() (resp *types.TestServiceResponse, err error) {  // Call the next service and get the response
    in_header := h.Ctx.RequestData.Request.Header  // The headers of the incomming requests
    
    ho, p, err := h.FindNextService();
    if err != nil {
        return nil, err
    }
    ep := strings.Join([]string{ho, p}, ":")
    req_url := strings.Join([]string{"http:/", ep, "api", h.Ctx.SvcOther}, "/")
    req, err := http.NewRequest("GET", req_url, nil); 
    if err != nil {
        return nil, err
    }
    err = PropTraceInfo(&in_header, &req.Header)
    if err != nil {
        return nil, err
    }

    client := &http.Client{}
    http_resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer http_resp.Body.Close()
    resp_body, err := ioutil.ReadAll(http_resp.Body)
    if err != nil {
        return nil, err
    }
    return types.RespDecode(resp_body)    
}