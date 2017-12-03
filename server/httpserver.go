package main

import (
    "fmt"
    "net/http"
    //"strings"
    "strconv"
    "bytes"
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
                  account.Name   = t.Account.Name
                  account.Input_goal= t.Account.Input_goal
                  account.Output_goal= t.Account.Output_goal
                  fmt.Println("intput goal"+account.Gender)
                  fmt.Println("intput goal"+t.Account.Output_goal)
                  fmt.Println("intput goal"+account.Output_goal)      
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
                fiti, _ := strconv.Atoi(t.Fitness.Input)
                nfiti,_ := strconv.Atoi(new_fitness.Input)
                newinput :=  strconv.Itoa(fiti+nfiti)
                fito, _ := strconv.Atoi(t.Fitness.Output)
                nfito,_ := strconv.Atoi(new_fitness.Output)
                newoutput :=  strconv.Itoa(fito+nfito)
                db.Model(&new_fitness).Update("input", newinput)
                db.Model(&new_fitness).Update("output", newoutput)
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
    }else if(t.Act=="CA"){
        CheckDb(t.UserID,&reply,db)             
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
    }else if(t.Act=="APIF"){
        client := &http.Client{}
        resp, _:= http.NewRequest("GET","https://trackapi.nutritionix.com/v2/search/instant?query="+t.Nutrition,nil)
        resp.Header.Add("Content-Type","application/json")
        resp.Header.Add("x-app-id","11c36a20")
        resp.Header.Add("x-app-key","bde581bab71e8481c6656ece20287122")
        r, _ := client.Do(resp)
        data, _ := ioutil.ReadAll(r.Body)
        w.Write(data)
        
        // calorie kcal
    }else if(t.Act=="APIE"){
        client := &http.Client{}
        var req Apie 
        var account Account 
        if errci:= db.Where("user_id = ?", t.UserID).First(&account).Error; errci==nil{
        
        req.Query = t.Exercise
        req.Gender = account.Gender
        req.Age= account.Age
        fmt.Println(account.Weight)
        req.Weight_kg,_= strconv.Atoi(account.Weight)
        req.Height_cm,_= strconv.Atoi(account.Height)
        
        res1D := &req
        res1B, _ := json.Marshal(res1D)
        fmt.Println(string(res1B))
        body := bytes.NewReader(res1B)
        resp, _:= http.NewRequest("POST","https://trackapi.nutritionix.com/v2/natural/exercise",body)
        resp.Header.Add("Content-Type","application/json")
        resp.Header.Add("x-app-id","11c36a20")
        resp.Header.Add("x-app-key","bde581bab71e8481c6656ece20287122")
        r, _ := client.Do(resp)
        data, _ := ioutil.ReadAll(r.Body)
        w.Write(data)
                                                           
        }else{
              reply.Error= "User doesn't exist"
        } 
    }else{
            reply.Error="Bad Request"
            // CheckDb(t.UserID,&reply)   
            // CheckFriend(t.UserID,&reply)   
            // CheckFitness(t.UserID,&reply)     
    }            

    res1D := &reply
    output,_ := json.Marshal(res1D)
    if(t.Act!="APIF"&&t.Act!="APIE"){
       w.Write(output)
    }
    if(t.Act=="APIE"&& reply.Error!=""){
       w.Write(output)
    }
    

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
        fitnesslist = append(fitnesslist,fitness[i].Date+" "+"Input "+fitness[i].Input+" Output "+fitness[i].Output)
      }      
      reply.Fitnesslist= fitnesslist
}
func  CheckDb( UserID string,reply *Reply, db *gorm.DB) {
      var accountinfo Account
      if err:=db.Where("user_id = ?", UserID).First(&accountinfo).Error; err!=nil{
        reply.Error = "no user match!"
      }
      reply.UserID = accountinfo.UserID
      reply.Password = accountinfo.Password
      reply.Height = accountinfo.Height
      reply.Weight = accountinfo.Weight
      reply.Gender = accountinfo.Gender
      reply.Age = accountinfo.Age
      reply.Name = accountinfo.Name
      reply.Input_goal = accountinfo.Input_goal
      reply.Output_goal = accountinfo.Output_goal

}

func prettyprint(b []byte) ([]byte, error) {
    var out bytes.Buffer
    err := json.Indent(&out, b, "", "  ")
    return out.Bytes(), err
}

func main() {
    var err error
    
    db,err := gorm.Open("mysql", "Healthpetbackup:Healthpetbackup@(healthpetbackup.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpetbackup?charset=utf8&parseTime=True&loc=Local")
    
    // fmt.Println(err)
    db.AutoMigrate(&Fitness{},&Account{},&Friend{})
    
    http.HandleFunc("/", Server) // set router
    err = http.ListenAndServe(":9090", nil) // set listen port
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
    Nutrition string 
    Exercise string 
    
}

type Reply struct{  
   UserID string `json:"UserID"`
   Name string 
   Password string  `json:"Password"`
   Height string `json:"Height"`
   Weight string `json:"Weight"`
   Gender string `json:"Gender"`
   Age  string `json:"Age"`
   Fitnesslist []string `json:"Fitnesslist"`
   Friendlist []string `json:"Friendlist"`
   Input_goal string
   Output_goal string
   Error string   
   
}
type Account struct {
    
    UserID string  `gorm:"primary_key"`
    Password string
    Height string 
    Weight string 
    Gender string 
    Age  string
    Name string
    Input_goal string
    Output_goal string
}
type Friend struct {
    
    UserID string  `gorm:"primary_key"`
    FriedID string `gorm:"primary_key"` 
}
type Fitness struct {
   
    Date string    `gorm:"primary_key"`
    UserID string  `gorm:"primary_key"`
    Input string
    Output string 
}
type Api struct {
   Query string `json:"query"`
   Timezone string `json:"timezone"`
}

type Apie struct {
  Query string `json:"query"`
  Gender string `json:"gender"`
  Weight_kg int `json:"weight_kg"`
  Height_cm int `json:"height_cm"`
  Age string       `json:"age"`

}
// type NfromApi struct{
//   Foods []Nutrition `json:"foods"`
// }

// type Nutrition struct {
//    Food_name string `json:"food_name"`
//    Serving_qty string `json:"serving_qty"`
//    Serving_unit string `json:"serving_unit"`
//    Serving_weight_grams string `json:"serving_weight_grams"`
//    Calories_kcal int `json:"nf_calories"`
//    Toalfat_g int `json:"nf_total_fat"`
//    Saturatedfat_g int `json:"nf_saturated_fat"`
//    Cholesterol_g int   `json:"nf_cholesterol"`
//    Sodium_g int `json:"nf_sodium"`
//    Carbohydrate_g int `json:"nf_total_carbohydrate"`
//    Fiber_g int  `json:"nf_dietary_fiber"`
//    Sugars_g int  `json:"nf_sugars"`
//    Protein_g int `json:"nf_protein"`
//    Potassium_g int `json:"nf_potassium"`
// }

