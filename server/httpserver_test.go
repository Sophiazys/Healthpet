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
    
    requestLI := request("LI","123123","123123",account,fitness,friendlist)
    expected := `{"UserID":"123123","Password":"123123","Height":29,"Weight":30,"Gender":"Female","Age":18,"Fitnesslist":["2018-11-08 0"],"Friendlist":null,"Error":""}`
    testfunc( requestLI, "LI", expected,t)

    requestLIerr := request("LI","456456","123123",account,fitness,friendlist)
    expected = `{"UserID":"","Password":"","Height":0,"Weight":0,"Gender":"","Age":0,"Fitnesslist":null,"Friendlist":null,"Error":"wrong Password or UserID"}`
    testfunc( requestLIerr, "LIerr", expected,t)

    account.UserID="987665"
    requestCI := request("CI","987665","123123",account,fitness,friendlist)
    expected = `{"UserID":"987665","Password":"123123","Height":12,"Weight":12,"Gender":"","Age":0,"Fitnesslist":null,"Friendlist":null,"Error":""}`
    testfunc( requestCI, "CI", expected,t)

    requestCIerr := request("CI","999999","123123",account,fitness,friendlist)
    expected = `{"UserID":"","Password":"","Height":0,"Weight":0,"Gender":"","Age":0,"Fitnesslist":null,"Friendlist":null,"Error":"User doesn't exist"}`
    testfunc( requestCIerr, "CIerr", expected,t)

    requestAF := request("AF","123123","123123",account,fitness,friendlist)
    expected = `{"UserID":"123123","Password":"123123","Height":29,"Weight":30,"Gender":"Female","Age":18,"Fitnesslist":["2018-11-08 0"],"Friendlist":null,"Error":""}`
    testfunc( requestAF, "AF", expected,t)

    requestAFerr := request("AF","999999","123123",account,fitness,friendlist)
    expected = `{"UserID":"","Password":"","Height":0,"Weight":0,"Gender":"","Age":0,"Fitnesslist":null,"Friendlist":null,"Error":"User doesn't exist"}`
    testfunc( requestAFerr, "AFerr", expected,t)

    requestFOerr := request("FO","999999","123123",account,fitness,friendlist)
    expected = `{"UserID":"","Password":"","Height":0,"Weight":0,"Gender":"","Age":0,"Fitnesslist":null,"Friendlist":null,"Error":"User doesn't exist"}`
    testfunc( requestFOerr, "FOerr", expected,t)










    
    
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

func tojson(request React_request) *bytes.Buffer {
    res1D := &request
    output,_ := json.Marshal(res1D)
    b := bytes.NewBuffer(output)
    return b
}

func testfunc(request React_request, Act string, expected string,t *testing.T){
    req, err := http.NewRequest("POST", "/", tojson(request))
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
    if rr.Body.String() != expected {
         t.Errorf( Act +  " returned unexpected body: got %v want %v",
             rr.Body.String(), expected)
    }else{
        fmt.Println("\n"+Act+" test")
    }
}




