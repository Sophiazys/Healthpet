package main

import (
    "fmt"
    "net/http"
    //"strings"
    //"strconv"
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
    var reply_friend Reply_friend
    var reply_fitness Reply_fitness

    if(t.Act=="LI"){            
            var account Account 
            if errci:=db.Where("user_id = ? AND password = ?", t.UserID,t.Password).First(&account).Error;errci==nil{
                reply= CheckDb(t.UserID)
                reply_friend=CheckFriend(t.UserID)
                reply_fitness=CheckFitness(t.UserID)
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
                reply_friend=CheckFriend(t.UserID)
                reply_fitness=CheckFitness(t.UserID)
            }            
    }else{
            reply=CheckDb(t.UserID)   
            reply_friend=CheckFriend(t.UserID)   
            reply_fitness=CheckFitness(t.UserID)     
    }            
    
    fmt.Println("before marshal to json")
    output,_ := json.Marshal(reply)
    w.Write(output)

    output2,_ := json.Marshal(reply_friend.friends)
    w.Write(output2)
    output3,_ := json.Marshal(reply_fitness.fitness)
    w.Write(output3)

}
func  CheckFriend( UserID string) Reply_friend{
    var reply_friend Reply_friend
    var friend  []Friend
    var friendslist string
    if err:= db.Where("user_id = ?", UserID).Find(&friend).Error; err!=nil{
         fmt.Println("friend check err")
      }
      for i:=range friend {
        fmt.Println("pengyou"+friend[i].UserID+" "+friend[i].FriedID)
        friendslist=friendslist + friend[i].FriedID+","
      }
      reply_friend.friends=friendslist

      return reply_friend
}

func  CheckFitness( UserID string) Reply_fitness{
    var reply_fitness Reply_fitness
    var fitness  []Fitness
    var fitnesslist string
    if err:= db.Where("user_id = ?", UserID).Find(&fitness).Error; err!=nil{
         fmt.Println("fitness check err")
      }
      for i:=range fitness {
        fitnesslist=fitnesslist + fitness[i].Date+" "+fitness[i].Calorie+","
      }
      reply_fitness.fitness=fitnesslist

      return reply_fitness
}
func  CheckDb( UserID string) Reply{
      var reply Reply
      var accountinfo Account
      if err:=db.Where("user_id = ?", UserID).First(&accountinfo).Error; err!=nil{
        fmt.Println("no user match!")
      }
    reply.UserID = accountinfo.UserID
      reply.Password = accountinfo.Password
      reply.Height = accountinfo.Height
      reply.Weight = accountinfo.Weight
      reply.Gender = accountinfo.Gender
      reply.Age = accountinfo.Age
      
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
}

type Reply_friend struct{  
   friends  string 
}

type Reply_fitness struct{  
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
    Calorie string
}
