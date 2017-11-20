// The response body structure used between test services

package types

import (
    "errors"
    "encoding/json"
    "fmt"
    "strconv"
    
    "utils"
    
    "github.com/tidwall/gjson"
)

type ServiceStatus struct {  // The service status of the local service
    Order string `json:"order"`  // The order of the local service from the end of the service chain
    ServName string `json:"service"`  // The name of the local service
    Version string `json:"version"`  // The version of the local service
}

type TestServiceResponse struct {   // The response body structure used between test services
    Len string `json:"length"`  // The length of the test services from the local service to the end of the service chain
    Chain []*ServiceStatus `json:"chain"` // The status of the services from the local service to the end of the service chain
}

func RespEncode(r *TestServiceResponse) (b []byte, err error) {  // Encode a response structure to a response body
    return json.Marshal(r)
}

func RespDecode(b []byte) (r *TestServiceResponse, err error) {  // Decode a response body to a response structure
    /*
    var dat map[string]interface{}
    
    if err := json.Unmarshal(b, &dat); err != nil {
        return nil, err
    }
    
    l = dat["length"].(int)
    c = dat["chain"].([]interface{})
    cha := make([]*ServiceStatus, len(c))
    for i, s := range c {
        if st, ok := s.(ServiceStatus); ok {
            cha[i] = &st
        }
    }
    
    return &TestServiceResponse{l, cha}, nil
    */
    
    json_str := utils.Convert(b)
    fmt.Println("Get a response: ", json_str)
    l := gjson.Get(json_str, "length").String()    
    cha_r := gjson.Get(json_str, "chain").Array()
    
    l_int , err := strconv.Atoi(l)
    if err != nil {
        return nil, err
    }
    if l_int != len(cha_r) {
       return nil, errors.New("The length of the service chain is not matched!")
    }
    
    cha := make([]*ServiceStatus, len(cha_r))
    for i, st_r := range cha_r {
        cha[i] = &ServiceStatus{gjson.Get(st_r.String(), "order").String(), 
                                gjson.Get(st_r.String(), "service").String(), 
                                gjson.Get(st_r.String(), "version").String()}
    }
    
    return &TestServiceResponse{l, cha}, nil   
}