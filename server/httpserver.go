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
// var db *gorm.DB

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func Server(w http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method)
    
    body, _ := ioutil.ReadAll(req.Body)
    var t React_request
    json.Unmarshal(body, &t)
    fmt.Println(t.Act)
    fmt.Println(t)
    var reply Reply
    db,_:= gorm.Open("mysql", "Healthpetbackup:Healthpetbackup@(healthpetbackup.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpetbackup?charset=utf8&parseTime=True&loc=Local")
    

    if(t.Act=="LI"){  //update log in info          
            var account Account 
            if errci:=db.Where("user_id = ? AND password = ?", t.UserID,t.Password).First(&account).Error;errci==nil{
                CheckDb(t.UserID,&reply,db)
                CheckFriend(t.UserID,&reply,db)
                CheckFitness(t.UserID,&reply,db)
            } else{
              reply.Error= "wrong Password or UserID"
            } 

    }else if(t.Act=="CI"){  //update account info
            var account Account        
            if errci:= db.Where("user_id = ?", t.UserID).First(&account).Error; errci==nil{
                if(account.UserID==t.Account.UserID){
                  account.Password = t.Account.Password
                  account.Height = t.Account.Height
                  account.Weight = t.Account.Weight
                  account.Gender = t.Account.Gender
                  account.Age    = t.Account.Age    
                  db.Save(&account) 
                }                            
                CheckDb(t.UserID,&reply,db) 
                CheckFriend(t.UserID,&reply,db)
                CheckFitness(t.UserID,&reply,db)
            }else{
              reply.Error= "User doesn't exist"
            } 

    }else if(t.Act=="AF"){  // update fitness info
            var new_fitness Fitness
            var account Account 
            if errci:= db.Where("user_id = ?", t.UserID).First(&account).Error; errci!=nil{
               reply.Error= "User doesn't exist"
            }else{
              if errd:= db.Where("user_id = ? AND date = ?",t.UserID,t.Fitness.Date).First(&new_fitness).Error; errd ==nil{
                fmt.Println(new_fitness)
                db.Model(&new_fitness).Update("calorie", t.Fitness.Calorie)
                //db.Save(&new_fitness)
              }else{
                fmt.Println(t.Fitness.Date)
                db.NewRecord(&t.Fitness)
                db.Create(&t.Fitness)
              }
              CheckDb(t.UserID,&reply,db) 
              CheckFriend(t.UserID,&reply,db)
              CheckFitness(t.UserID,&reply,db)
            }
            
    }else if(t.Act=="FO"){
            var account Account 
            if errci:= db.Where("user_id = ?", t.UserID).First(&account).Error; errci!=nil{
               reply.Error= "User doesn't exist"
            }else{
              for i:=range t.Friendlist {
                var tmpfriend Friend
                var checkfrend Friend 
                var faccount Account
                if errf:= db.Where("user_id = ?",t.Friendlist[i]).First(&faccount).Error; errf!=nil{
                  reply.Error = reply.Error+" friend "+ t.Friendlist[i]+ " does not exist"
                }else if errf:= db.Where("user_id = ? AND fried_id = ? ",t.UserID,t.Friendlist[i]).First(&checkfrend).Error; errf==nil{
                  //friend already exist need to add 
                }else{
                  tmpfriend.UserID=t.UserID
                  tmpfriend.FriedID=t.Friendlist[i]
                  db.NewRecord(&tmpfriend)
                  db.Create(&tmpfriend)
                }                               
              }
              CheckDb(t.UserID,&reply,db) 
              CheckFriend(t.UserID,&reply,db)
              CheckFitness(t.UserID,&reply,db)
            }
    }else if(t.Act=="SI"){  //update account info
            var newaccount Account        
            if errci:= db.Where("user_id = ?", t.UserID).First(&newaccount).Error; errci!=nil{
                newaccount.UserID = t.UserID
                newaccount.Password= t.Password
                db.NewRecord(&newaccount)
                db.Create(&newaccount)
            }else{
              reply.Error= "User already exist"
            } 
    }else{
            reply.Error="Bad Request"
            // CheckDb(t.UserID,&reply)   
            // CheckFriend(t.UserID,&reply)   
            // CheckFitness(t.UserID,&reply)     
    }            

    res1D := &reply
    output,_ := json.Marshal(res1D)
    w.Write(output)

}
func  CheckFriend( UserID string, reply *Reply, db *gorm.DB) {
    var friend  []Friend
    if err:= db.Where("user_id = ?", UserID).Find(&friend).Error; err!=nil{
         fmt.Println("friend check err")
    }
      for i:=range friend {
        reply.Friendlist=append(reply.Friendlist,friend[i].FriedID) 
      }      
}

func  CheckFitness( UserID string, reply *Reply, db *gorm.DB) {
    var fitness  []Fitness
    var fitnesslist []string
    if err:= db.Where("user_id = ?", UserID).Find(&fitness).Error; err!=nil{
         fmt.Println("fitness check err")
      }
      for i:=range fitness {
        fitnesslist = append(fitnesslist,fitness[i].Date+" "+fitness[i].Calorie)
      }      
      reply.Fitnesslist= fitnesslist
}
func  CheckDb( UserID string,reply *Reply, db *gorm.DB) {
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

}
func main() {
    var err error
    
    // db,err = gorm.Open("mysql", "Healthpetbackup:Healthpetbackup@(healthpetbackup.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpetbackup?charset=utf8&parseTime=True&loc=Local")
    
    // fmt.Println(err)
    // db.AutoMigrate(&Fitness{},&Account{},&Friend{})
    
    http.HandleFunc("/", Server) // set router
    err = http.ListenAndServe(":9191", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

type React_request struct {
    Act      string 
    UserID   string  
    Password string 
    Account  Account 
    Fitness  Fitness
    Friendlist []string 
}

type Reply struct{  
   UserID string `json:"UserID"`
   Password string  `json:"Password"`
   Height int `json:"Height"`
   Weight int `json:"Weight"`
   Gender string `json:"Gender"`
   Age  int `json:"Age"`
   Fitnesslist []string `json:"Fitnesslist"`
   Friendlist []string `json:"Friendlist"`
   Error string   
   
}
type Account struct {
    
    UserID string  `gorm:"primary_key"`
    Password string
    Height int 
    Weight int 
    Gender string 
    Age  int 
}
type Friend struct {
    
    UserID string  `gorm:"primary_key"`
    FriedID string `gorm:"primary_key"` 
}
type Fitness struct {
   
    Date string    `gorm:"primary_key"`
    UserID string  `gorm:"primary_key"`
    Calorie string 
}
