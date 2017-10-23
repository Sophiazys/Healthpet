// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
    "strconv"

	"flag"
	"html/template"
    "fmt"
	"log"
	"net/http"
    "strings"
	"github.com/gorilla/websocket"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var db *gorm.DB
var upgrader = websocket.Upgrader{} // use default options

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func echo(w http.ResponseWriter, r *http.Request) {
    fmt.Println("in echo function, connect to the front end")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		


        mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message)

        msg := string(message[:])
        msglist:= strings.Split(msg, " ")
        
        /////////////////////////////////////////////////////////////////////////////////
        if(msglist[0]=="LI"){
            var account Account
            found:=false           
            if errlogin:= db.Where("user_id = ? AND password = ?", msglist[1],msglist[2]).First(&account).Error; errlogin==nil{
                found=true
            }            
            if(found){
                err = c.WriteMessage(mt, []byte("log in success!"))

            }else{
                err = c.WriteMessage(mt, []byte("No record!")) 
            }            
        
        }else if(msglist[0]=="CI"){

            var account Account
            found:=false           
            if errci:= db.Where("user_id = ?", msglist[1]).First(&account).Error; errci==nil{
                found=true
            }            
            if(found){
                account.Height,_=strconv.Atoi(msglist[2])
                account.Weight,_=strconv.Atoi(msglist[3])
                account.Gender=msglist[4]
                account.Age,_=strconv.Atoi(msglist[5])
                db.Save(&account)
                err = c.WriteMessage(mt, []byte("change success"))
            }else{
                err = c.WriteMessage(mt, []byte("No record!")) 
            }  

        }else if(msglist[0]=="IC"){
            var fitness  Fitness
            var fitnessupdate Fitness
            found:=false           
            if erric:= db.Where("user_id = ?", msglist[1]).First(&fitness).Error; erric==nil{
                found=true
            }            
            if(found){
                fitnessupdate.UserID=msglist[1]
                fitnessupdate.Date=msglist[3]
                fitnessupdate.Calorie,_=strconv.Atoi(msglist[2])
                db.Save(&fitnessupdate)
                err = c.WriteMessage(mt, []byte("change fitness success"))
            }else{
                err = c.WriteMessage(mt, []byte("No record!")) 
            }  

        }else if(msglist[0]=="RC"){
            var fitness  []Fitness
            var fitnesshistory string
            found:=false           
            
            if erric:= db.Where("user_id = ?", msglist[1]).Find(&fitness).Error; erric==nil{
                found=true
            }            
            if(found){
                for i:=range fitness{
                    fitnessrecord:=  strconv.Itoa(fitness[i].Calorie)+":"+fitness[i].Date+";"
                    fitnesshistory=fitnesshistory+" "+fitnessrecord
                }
                
                err = c.WriteMessage(mt, []byte(fitnesshistory))
            }else{
                err = c.WriteMessage(mt, []byte("No record!")) 
            }  

        }else if(msglist[0]=="AF"){
            var friend  []Friend
            var friendlist string
            found:=false           
            
            if erric:= db.Where("user_id = ?", msglist[1]).Find(&friend).Error; erric==nil{
                found=true
            }            
            if(found){
                for i:=range friend{
                    
                    friendlist=friendlist+" "+friend[i].FriedID
                }
                
                err = c.WriteMessage(mt, []byte(friendlist))
            }else{
                err = c.WriteMessage(mt, []byte("No record!")) 
            }    
        }else{

        }            
        /////////////////////////////////////////////////////////////////////////////

    }
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
    fmt.Println("server start, now open the db")
    //db, _ = gorm.Open("mysql", "root:84921699@(localhost:3306)/mistro1?charset=utf8&parseTime=True&loc=Local")
    
    db,_= gorm.Open("mysql", "Healthpetbackup:Healthpetbackup@(healthpetbackup.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpetbackup?charset=utf8&parseTime=True&loc=Local")
    // sophia1 := Account{UserID: "yz3083", Password:"yingshuangzheng",Height: 29,Weight :30,Gender:"Female",Age: 18}
    // db.NewRecord(sophia1)
    // db.Create(&sophia1)
    // db.NewRecord(sophia1)
    //db, _ = gorm.Open("mysql", "root:84921699@(160.39.140.131:3306)/mistro1?charset=utf8&parseTime=True&loc=Local")
    
    db.AutoMigrate(&Fitness{})
	
    flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
// func forminfo(input User) string {
//     reply:= input.Name +"\n"
//     return reply
// }

type Account struct {
    gorm.Model 
    UserID string `gorm:"UserID"`
    Password string `gorm:"Password"`
    //Name string
    Height int `gorm:"Height"`
    Weight int `gorm:"Weight"`
    Gender string `gorm:"Gender"`
    Age  int `gorm:"Age"`    
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

