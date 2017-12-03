package main

import (
	"encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "strconv"
    "fmt"
    "bytes"
    "github.com/jinzhu/gorm"
    _"github.com/jinzhu/gorm/dialects/mysql"
)

func TestServer(t *testing.T) {
    db,_:= gorm.Open("mysql", "Healthpetbackup:Healthpetbackup@(healthpetbackup.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpetbackup?charset=utf8&parseTime=True&loc=Local")
    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
    // pass 'nil' as the third parameter.
    var account Account
    account.UserID="123123"
    account.Password="123123"
    account.Height= "12"
    account.Weight= "12"
    account.Input_goal="120"
    account.Output_goal="240" 
    account.Name=""

    var fitness Fitness
    fitness.UserID="AFtest"
    fitness.Date="2018-11-08"
    fitness.Input="120"
    fitness.Output="240"
    
    friendlist:= []string{"111111","123456"}
    ///////record LItest
    requestLI := request("LI","LItest","123123",account,fitness,friendlist,"egg","run")
    expected := `{"UserID":"LItest","Name":"","Password":"123123","Height":"29","Weight":"30","Gender":"Female","Age":"18","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":""}`
    testfunc( requestLI, "LI", expected,t)

    requestLIerr := request("LI","456456","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"wrong Password or UserID"}`
    testfunc( requestLIerr, "LIerr", expected,t)

    ///////////record CItest
    account.UserID="CItest"
    requestCI := request("CI","CItest","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"CItest","Name":"","Password":"123123","Height":"12","Weight":"12","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"120","Output_goal":"240","Error":""}`
    testfunc( requestCI, "CI", expected,t)

    requestCIerr := request("CI","999999","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"User doesn't exist"}`
    testfunc( requestCIerr, "CIerr", expected,t)

    ///////////record AFtest
    requestAF := request("AF","AFtest","123123",account,fitness,friendlist,"egg","run")
    expectAF(requestAF, "AF",t,db)
    
   
    ////////////////////////////////////////////
    requestAFerr := request("AF","999999","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"User doesn't exist"}`
    testfunc( requestAFerr, "AFerr", expected,t)

    /////////////record FOtest
    requestFO := request("FO","FOtest","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"FOtest","Name":"","Password":"123123","Height":"29","Weight":"30","Gender":"Female","Age":"18","Fitnesslist":null,"Friendlist":["111111"],"Input_goal":"","Output_goal":"","Error":" friend 123456 does not exist"}`
    testfunc( requestFO, "FO", expected,t)


    requestFOerr := request("FO","999999","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"User doesn't exist"}`
    testfunc( requestFOerr, "FOerr", expected,t)


    requestSI := request("SI","testsignin","123123",account,fitness,friendlist,"egg","run")
    req, err := http.NewRequest("POST", "/", tojson(requestSI))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Server)
    handler.ServeHTTP(rr, req)      
    var accountSI Account 
    if errSI:= db.Where("user_id = ?", "testsignin").Find(&accountSI).Error; errSI!=nil{
         fmt.Println("SI fail")
    }else{
        db.Delete(&accountSI)
        fmt.Println("SI OK")
    }
 
   
    

    /////////////// record sierrtest
    requestSIerr := request("SI","SIerrtest","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"User already exist"}`
    testfunc( requestSIerr, "SIerr", expected,t)

    requestActerr := request("Sophia","SIerrtest","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"Bad Request"}`
    testfunc( requestActerr, "Acterr", expected,t)

    ///////////////////test API///////////////
    requestAPIE := request("APIE","123123","123123",account,fitness,friendlist,"egg","run")
    testAPI( requestAPIE, "APIE", t)

    requestAPIEerr := request("APIE","APIEerr","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"","Name":"","Password":"","Height":"","Weight":"","Gender":"","Age":"","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":"User doesn't exist"}`
    testfunc( requestAPIEerr, "APIE",expected, t)

    ////////////////////////////////////////////////////////////////////////////////////////

    requestAPIF := request("APIF","APIFtest","123123",account,fitness,friendlist,"egg","run")
    testAPI( requestAPIF, "APIF", t)

    ///////////////////test CA///////////////
    requestCA := request("CA","123123","123123",account,fitness,friendlist,"egg","run")
    expected = `{"UserID":"123123","Name":"","Password":"123123","Height":"29","Weight":"30","Gender":"Female","Age":"18","Fitnesslist":null,"Friendlist":null,"Input_goal":"","Output_goal":"","Error":""}`
    testfunc( requestCA, "CA", expected,t)
    
}
func request(Act string, UserID string, Password string, account Account,fitness Fitness,friendlist []string, Nutrition string, Exercise string) React_request{
    var r React_request
    r.Act= Act
    r.UserID = UserID
    r.Password= Password
    r.Account = account
    r.Account.UserID= UserID
    r.Fitness = fitness
    r.Friendlist=friendlist
    r.Nutrition=Nutrition
    r.Exercise = Exercise
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

func testAPI(request React_request, Act string, t *testing.T){
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
    if len(rr.Body.String())==0 {
         t.Errorf( Act +  "fail, no result return")
    }else{
        fmt.Println("\n"+Act+" test")
    }
}

func expectAF(request React_request, Act string, t *testing.T, db *gorm.DB) {
    
    var fitness Fitness
    var fitnessafter Fitness
    db.Where("user_id = ? AND date = ?","AFtest","2018-11-08").First(&fitness); 

    i, _ :=strconv.Atoi(fitness.Input)
    o, _ :=strconv.Atoi(fitness.Output)
    i=i+120
    o=o+240

     req, err := http.NewRequest("POST", "/", tojson(request))
    if err != nil {
        t.Fatal(err)
    }    
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Server)
    handler.ServeHTTP(rr, req)

    db.Where("user_id = ? AND date = ?","AFtest","2018-11-08").First(&fitnessafter)
   
    if strconv.Itoa(i)==fitnessafter.Input && strconv.Itoa(o)==fitnessafter.Output{
         fmt.Println("\n"+Act+" test")
    }else{
        t.Errorf( Act +  "fail, fail to update input/ output calorie")
    }
}




