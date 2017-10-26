package main

import (
    "fmt"
    "net/http"
    //"strings"
    "strconv"
    "log"
    "encoding/json"
    "io/ioutil"
    "github.com/jinzhu/gorm"
    _"github.com/jinzhu/gorm/dialects/mysql"

)
var db *gorm.DB

func server(w http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method)
    
    body, _ := ioutil.ReadAll(req.Body)
    var t React_request
    json.Unmarshal(body, &t)
    fmt.Println(t.Act)
    fmt.Println(t.Password)
    var reply Reply

    if(t.Act=="LI"){            
            var account Account 
            if errci:=db.Where("user_id = ? AND password = ?", t.UserID,t.Password).First(&account).Error;errci==nil{
                reply= CheckDb(t.UserID)
            }           
            fmt.Println("after query")

    }else if(t.Act=="CI"){
            var account Account        
            if errci:= db.Where("user_id = ?", t.UserID).First(&account).Error; errci==nil{
                account = t.account
                db.Save(&account)
                fmt.Println("found" + t.UserID)
                fmt.Println("found acc"+account.UserID)
                reply = CheckDb(t.UserID) 
            }            
    }else{
                reply=CheckDb(t.UserID)
            
    }            
    
    fmt.Println("before marshal to json")
    output,_ := json.Marshal(reply)
    w.Write(output)
}

func  CheckDb( UserID string) Reply{
  var reply Reply
  var friend  []Friend
  var accountinfo Account
  var friendslist string
  var fitnesslist string
  if err:=db.Where("user_id = ?", UserID).First(&accountinfo).Error; err!=nil{
    fmt.Println("no user match!")
  }

  reply.UserID = accountinfo.UserID
  reply.Password = accountinfo.Password
  reply.Height = accountinfo.Height
  reply.Weight = accountinfo.Weight
  reply.Gender = accountinfo.Gender
  reply.Age = accountinfo.Age
  
  if err:= db.Where("user_id = ?", UserID).Find(&friend).Error; err!=nil{
     fmt.Println("friend check err")
  }
  for i:=range friend {
    friendslist=friendslist + friend[i].FriedID+","
  }

  var fitness  []Fitness                         
  if err := db.Where("user_id = ?", UserID).Find(&fitness).Error; err!=nil{
      fmt.Println("friend check err")
  } 
  for i:=range fitness{
     fitnesslist = fitnesslist+ " "+ fitness[i].Date+":"+strconv.Itoa(fitness[i].Calorie)
  }
  reply.friends=friendslist
  reply.fitness=fitnesslist
  
  return reply
}

func main() {
    var err error
    
    db,err = gorm.Open("mysql", "Healthpetbackup:Healthpetbackup@(healthpetbackup.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpetbackup?charset=utf8&parseTime=True&loc=Local")
    
    fmt.Println(err)
    db.AutoMigrate(&Fitness{})
    
    http.HandleFunc("/", server) // set router
    err = http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

type React_request struct {
    Act      string 
    UserID   string  
    Password string 
    account   Account
    
}

type Reply struct{  
   UserID string 
   Password string
   Height int 
   Weight int 
   Gender string 
   Age  int 
   friends  string 
   fitness  string

}

type Account struct {
    gorm.Model 
    UserID string 
    Password string
    Height int 
    Weight int 
    Gender string 
    Age  int 
}
type Friend struct {
    gorm.Model 
    UserID string 
    FriedID string     
}
type Fitness struct {
    gorm.Model 
    Date string
    UserID string 
    Calorie int   
}
