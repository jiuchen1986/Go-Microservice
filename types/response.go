// The response body structure used between test services

package types

import (
//    "errors"
    "encoding/json"
    "fmt"

    "utils"    
//    "github.com/jiuchen1986/Go-Microservice/utils"    
    "github.com/tidwall/gjson"
)

type ServiceStatus struct {  // The service status of the local service
    Order string `json:"order"`  // The order of the local service from the end of the service chain
    ServName string `json:"service"`  // The name of the local service
    Version string `json:"version"`  // The version of the local service
}

type ServiceChain struct {  // A service chain
     Starter string `json:"starter"`  // The name of the starting service of this chain
     ChainId string `json:"id"` // A chain id unique in a service
     Chain []*ServiceStatus `json:"chain"` // The status of the services from the local service to the end of the service chain
     Len string `json:"length"`  // The length of the test services from the local service to the end of the service chain
}

type TestServiceResponse struct {   // The response body structure used between test services
    MainChain *ServiceChain `json:"main_chain"`  // The main chain
    SubChains []*ServiceChain `json:"sub_chains"` // The sub chains involved in this call
}

func RespEncode(r *TestServiceResponse) (b []byte, err error) {  // Encode a response structure to a response body
    return json.Marshal(r)
}

func RespDecode(b []byte) (r *TestServiceResponse, err error) {  // Decode a response body to a response structure    
    json_str := utils.Convert(b)
    // fmt.Println("response.RespDecode: Get a response: ", json_str)    
    
    main_chain_gjson_str := gjson.Get(json_str, "main_chain").String()
    main_chain_resp := &ServiceChain{gjson.Get(main_chain_gjson_str, "starter").String(), 
                                     gjson.Get(main_chain_gjson_str, "id").String(), 
                                     make([]*ServiceStatus, 0), 
                                     gjson.Get(main_chain_gjson_str, "length").String()}
    var status_gjson_str string                                     
    for _, st_gjson := range gjson.Get(main_chain_gjson_str, "chain").Array() {
        status_gjson_str = st_gjson.String()
        main_chain_resp.Chain = append(main_chain_resp.Chain, &ServiceStatus{gjson.Get(status_gjson_str, "order").String(), 
                                                                             gjson.Get(status_gjson_str, "service").String(), 
                                                                             gjson.Get(status_gjson_str, "version").String()}) 
    }
    
    sub_chains_gjson := gjson.Get(json_str, "sub_chains").Array()
    sub_chains_resp := make([]*ServiceChain, 0)
    
    // temp var for the chains loop
    var chain_resp []*ServiceStatus
    var chain_gjson_str string   
    
    for _, chain_gjson := range sub_chains_gjson {
        chain_gjson_str = chain_gjson.String()        
        chain_resp = make([]*ServiceStatus, 0)
        
        for _, status_gjson := range gjson.Get(chain_gjson_str, "chain").Array() {
            status_gjson_str = status_gjson.String()
            chain_resp = append(chain_resp, &ServiceStatus{gjson.Get(status_gjson_str, "order").String(), 
                                              gjson.Get(status_gjson_str, "service").String(), 
                                              gjson.Get(status_gjson_str, "version").String()})
        }
        sub_chains_resp = append(sub_chains_resp, &ServiceChain{gjson.Get(chain_gjson_str, "starter").String(), 
                                          gjson.Get(chain_gjson_str, "id").String(), 
                                          chain_resp,  
                                          gjson.Get(chain_gjson_str, "length").String()})
    }
    
    fmt.Println("response.RespDecode: Successfully decode the response: ", json_str)
    return &TestServiceResponse{main_chain_resp, sub_chains_resp}, nil   
}