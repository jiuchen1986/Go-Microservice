// Defining the handler processing the requests for the service chain

package handler

import (
    "fmt"
    "strings"
    "net/http"
    "io/ioutil"
    "strconv"
    "time"
    "math/rand"
    
    "app"
    "types"
    "utils"
//    "github.com/jiuchen1986/Go-Microservice/app"
//    "github.com/jiuchen1986/Go-Microservice/types"
//    "github.com/jiuchen1986/Go-Microservice/utils"
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

func (h *ServiceChainHandler) Process(delay time.Duration) error {  // the main requests process of the handler
    fmt.Println("handler.servchainhandler: Delay for ", delay)
    time.Sleep(delay)
    
    v, er := h.VerifyPath()
    if er != nil {
        h.Ctx.NotFound()
        return er
    }
    if !v {
        fmt.Println("handler.servchainhandler: Invailid incomming path")
        return h.Ctx.NotFound()
    }
    
    var main_resp, sub_resp *types.TestServiceResponse
    ch_resp := make([]chan *types.TestServiceResponse, 2)
    
    ch_resp[0] = make(chan *types.TestServiceResponse)
    go h.FollowMainChain(ch_resp[0])
    ch_resp[1] = make(chan *types.TestServiceResponse)
    go h.FollowSubChain(ch_resp[1])
    main_resp, sub_resp = <-ch_resp[0], <-ch_resp[1]
    
    if main_resp == nil {
        fmt.Println("handler.servchainhandler: No main response")
        return h.Ctx.NotFound()
    }
    
    // resp_b, _ := types.RespEncode(main_resp)
    // fmt.Printf("handler.servchainhandler: get the main response: %s\n", utils.Convert(resp_b))
    
    stat_lo, er := GetLocalServiceStatus()
    if er != nil {
        h.Ctx.NotFound()
        return er
    }
    
    rand_factor := (int)(stat_lo.ServName[len(stat_lo.ServName) - 1]) // using the ascii code of the last letter 
                                                                      // of the local serivce name as a rand factor
    rand_factor_main := (int)("main"[len("main") -1])
    rand_factor_sub := (int)("sub"[len("sub") -1])
    
    // process the response from the main chain
    chain_id, er := strconv.Atoi(main_resp.MainChain.ChainId)
    if er != nil {
        h.Ctx.NotFound()
        return er
    }
    rand.Seed((int64)(chain_id + rand_factor + rand_factor_main))
    main_resp.MainChain.ChainId = strconv.Itoa(rand.Int())
    main_resp.MainChain.Starter = stat_lo.ServName
    chain_l, er := strconv.Atoi(main_resp.MainChain.Len)
    if er != nil {
        h.Ctx.NotFound()
        return er
    }
    stat_lo.Order = strconv.Itoa(chain_l + 1)
    main_resp.MainChain.Len = stat_lo.Order
    main_resp.MainChain.Chain = append(main_resp.MainChain.Chain, stat_lo)
    
    // process the response from the sub chain    
    if sub_resp != nil {        
        
        // process the main chain in the response from the sub chain
        stat_lo, _ = GetLocalServiceStatus()
        chain_id, er := strconv.Atoi(sub_resp.MainChain.ChainId)
        if er != nil {
            h.Ctx.NotFound()
            return er
        }
        rand.Seed((int64)(chain_id + rand_factor + rand_factor_sub))
        sub_resp.MainChain.ChainId = strconv.Itoa(rand.Int())
        sub_resp.MainChain.Starter = stat_lo.ServName
        chain_l, er := strconv.Atoi(sub_resp.MainChain.Len)
        if er != nil {
            h.Ctx.NotFound()
            return er
        }
        stat_lo.Order = strconv.Itoa(chain_l + 1)
        sub_resp.MainChain.Len = stat_lo.Order
        sub_resp.MainChain.Chain = append(sub_resp.MainChain.Chain, stat_lo)
        main_resp.SubChains = append(main_resp.SubChains, sub_resp.MainChain)
        
        // process the sub chain in the response from the sub chain
        for _, sub_chain := range sub_resp.SubChains {
            stat_lo, _ = GetLocalServiceStatus()
            chain_id, er = strconv.Atoi(sub_chain.ChainId)
            if er != nil {
                h.Ctx.NotFound()
                return er
            }
            rand.Seed((int64)(chain_id + rand_factor + rand_factor_sub))
            sub_chain.ChainId = strconv.Itoa(rand.Int())
            sub_chain.Starter = stat_lo.ServName
            chain_l, er = strconv.Atoi(sub_chain.Len)
            if er != nil {
                h.Ctx.NotFound()
                return er
            }
            stat_lo.Order = strconv.Itoa(chain_l + 1)
            sub_chain.Len = stat_lo.Order
            sub_chain.Chain = append(sub_chain.Chain, stat_lo)
            main_resp.SubChains = append(main_resp.SubChains, sub_chain) 
        }
    }

    resp_b, er := types.RespEncode(main_resp)
    if er != nil {
        h.Ctx.NotFound()
        return er
    } else {
            
        fmt.Printf("handler.servchainhandler: Send a response with OK!\n")
        fmt.Printf("handler.servchainhandler: Response body: \n")
        fmt.Println(utils.Convert(resp_b))
            
        return h.Ctx.OK(resp_b)
    }

    return nil    
}

func (h *ServiceChainHandler) VerifyPath() (v bool, err error) {  // Verify the incomming path
    st_lo, err := GetLocalServiceStatus()
    if err != nil {
        return false, err
    }
    if strings.Compare(strings.Split(h.Ctx.SvcLo, "_")[0], st_lo.ServName) != 0 {
        return false, nil
    }
    return true, nil        
}

func (h *ServiceChainHandler) FindNextServiceMain() (url string, err error) {  // Return the url for the next service 
                                                                               // in the main chain
    svcTo := strings.Split(strings.Split(h.Ctx.SvcOther, "/")[0], "_")[0]
    fmt.Println("handler.servchainhandler: The next service in the main chain: ", svcTo)
    /*
    return strings.Join([]string{"http:/", 
                                 strings.Join([]string{svcTo, "8082"}, ":"), 
                                 "api", 
                                 h.Ctx.SvcOther}, "/"), nil
    */
    
    return "http://10.0.2.15:8082/api/" + h.Ctx.SvcOther, nil
}

func (h *ServiceChainHandler) FindNextServiceSub() (url string, err error) {  // Return the url for the next service 
                                                                              // in the sub chain
    sub_chain := strings.Split(h.Ctx.SvcLo, "_")
    if len(sub_chain) < 2 {
        return "", nil
    }
    fmt.Println("handler.servchainhandler: The next service in the sub chain: ", sub_chain[1])
    /*
    return strings.Join([]string{"http:/", 
                                 strings.Join([]string{sub_chain[1], "8082"}, ":"), 
                                 "api", 
                                 strings.Join(sub_chain[1:], "/")}, "/"), nil
    */
    
    return "http://10.0.2.15:8082/api/" + strings.Join(sub_chain[1:], "/"), nil
}

func PropTraceInfo(ih, oh *http.Header) error {  // Collect and progapate the headers from the incomming request to the outgoing request for tracing
    for _, h := range trace_header {
        if v, ok := (*ih)[h]; ok {
            // fmt.Println("Found an incomming header: ", h, "=", v[0])
            if v[0] != "" {
                (*oh).Add(h, v[0])
                // fmt.Println("Add a header: ", h, "=", v[0])
            }
        }
    }
    return nil
}

func (h *ServiceChainHandler) FollowMainChain(ch chan *types.TestServiceResponse) error {  // Call the next service and get the response
                                                                                           // in the main chain
    
    in_header := h.Ctx.RequestData.Request.Header  // The headers of the incomming requests
    var resp *types.TestServiceResponse = nil
    defer func(){ ch <- resp }()
    
    req_url, err := h.FindNextServiceMain();
    if err != nil {
        return err
    }
    fmt.Println("handler.servchainhandler: The request url for the main chain: ", req_url)
    req, err := http.NewRequest("GET", req_url, nil); 
    if err != nil {
        return err
    }
    err = PropTraceInfo(&in_header, &req.Header)
    if err != nil {
        return err
    }

    client := &http.Client{}
    http_resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer http_resp.Body.Close()
    resp_body, err := ioutil.ReadAll(http_resp.Body)
    if err != nil {
        return err
    }
    resp, err = types.RespDecode(resp_body)
    if err != nil {
        resp = nil
        return err
    }
    resp_b, _ := types.RespEncode(resp)
    fmt.Printf("handler.servchainhandler: get the main response: %s\n", utils.Convert(resp_b))
    return nil    
}

func (h *ServiceChainHandler) FollowSubChain(ch chan *types.TestServiceResponse) error {  // Call the next service and get the response
                                                                                          // in the sub chain
    
    in_header := h.Ctx.RequestData.Request.Header  // The headers of the incomming requests
    var resp *types.TestServiceResponse = nil
    defer func(){ ch <- resp }()
    
    req_url, err := h.FindNextServiceSub();
    if err != nil {
        return err
    }
    if strings.Compare(req_url, "") == 0 {
        return nil
    }
    fmt.Println("handler.servchainhandler: The request url for the sub chain: ", req_url)
    req, err := http.NewRequest("GET", req_url, nil); 
    if err != nil {
        return err
    }
    err = PropTraceInfo(&in_header, &req.Header)
    if err != nil {
        return err
    }

    client := &http.Client{}
    http_resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer http_resp.Body.Close()
    resp_body, err := ioutil.ReadAll(http_resp.Body)
    if err != nil {
        return err
    }
    resp, err = types.RespDecode(resp_body)
    if err != nil {
        resp = nil
        return err
    }
    resp_b, _ := types.RespEncode(resp)
    fmt.Printf("handler.servchainhandler: get the sub response: %s\n", utils.Convert(resp_b))
    return nil    
}