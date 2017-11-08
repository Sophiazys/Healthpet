package main

import (
	"encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "fmt"
    "bytes"
)

func TestServer(t *testing.T) {
    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
    // pass 'nil' as the third parameter.
    var request_front React_request
    var account Account
    account.UserID="123123"
    account.Password="123123"
    account.Height= 12
    account.Weight= 12

    var fitness Fitness
    fitness.UserID="123123"
    fitness.Date="2018-11-08"
    fitness.Calorie="0"
    
    friendlist:= []string{"yz3083","123456"}
    
    request_front = request("LI","123123","123123",account,fitness,friendlist)
	res1D := &request_front
    output,_ := json.Marshal(res1D)
    b := bytes.NewBuffer(output)

    req, err := http.NewRequest("POST", "/", b)
    if err != nil {
        t.Fatal(err)
    }

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Server)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
    // directly and pass in our Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)
    // Check the response body is what we expect.
    expected := `{"UserID":"123123","Password":"123123","Height":29,"Weight":30,"Gender":"Female","Age":18,"Fitnesslist":null,"Friendlist":null,"Error":""}`
    if rr.Body.String() != expected {
         t.Errorf("handler returned unexpected body: got %v want %v",
             rr.Body.String(), expected)
    }else{
        fmt.Println("\nLI test")
    }
}
func request(Act string, UserID string, Password string, account Account,fitness Fitness,friendlist []string ) React_request{
    var r React_request
    r.Act= Act
    r.UserID = UserID
    r.Password= Password
    r.Account = account
    r.Fitness = fitness
    r.Friendlist=friendlist
    return r 
}



// type React_request struct {
//     Act      string 
//     UserID   string  
//     Password string 
//     Account  Account 
//     Fitness  Fitness
//     Friendlist []string 
// }

// type Fitness struct {   
//     Date string    
//     UserID string  
//     Calorie string 
// }

// type Account struct {   
//     UserID string  
//     Password string
//     Height int 
//     Weight int 
//     Gender string 
//     Age  int 
// }


