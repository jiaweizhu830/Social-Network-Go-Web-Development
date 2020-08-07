 1 package main
  2 
  3 import (
  4 >       "encoding/json"
  5 >       "fmt"
  6 >       "net/http"
  7 >       "reflect"
  8 >       "regexp"
  9 >       "time"
 10 
 11 >       jwt "github.com/dgrijalva/jwt-go"
 12 >       "github.com/olivere/elastic"
 13 )
 14 
 15 const (
 16 >       USER_INDEX = "user"
 17 )
 18 
 19 type User struct {
 20 >       Username string `json:"username"`
 21 >       Password string `json:"password"`
 22 >       Age      int64  `json:"age"`
 23 >       Gender   string `json:"gender"`
 24 }
 25 
 26 var mySigningKey = []byte("secret")
 27 
 28 func checkUser(username, password string) (bool, error) {
 29 >       query := elastic.NewTermQuery("username", username)
 30 >       searchResult, err := readFromES(query, USER_INDEX)
 31 >       if err != nil {
 32 >       >       return false, err
 33 >       }
 34 
 35 >       var utype User
 36 
 37 >       for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
 38 >       >       if u, ok := item.(User); ok {
 39 >       >       >       if u.Password == password {
 40 >       >       >       >       fmt.Printf("Login as %s\n", username)
 41 >       >       >       >       return true, nil
 42 >       >       >       }
 43 >       >       }
 44 >       }
 45 >       return false, nil
 46 }
 47 
 48 func addUser(user *User) (bool, error) {
 49 >       query := elastic.NewTermQuery("username", user.Username)
 50 >       searchResult, err := readFromES(query, USER_INDEX)
 51 >       if err != nil {
 52 >       >       return false, err
 53 >       }
 54 
 55 >       //check if user.Username is already in DB
 56 >       if searchResult.TotalHits() > 0 {
 57 >       >       return false, nil
 58 >       }
 59 
 60 >       err = saveToES(user, USER_INDEX, user.Username)
 61 >       if err != nil {
 62 >       >       return false, err
 63 >       }
 64 >       fmt.Printf("User is added: %s\n", user.Username)
 65 >       return true, nil
 66 }
 67 
 68 func handlerLogin(w http.ResponseWriter, r *http.Request) {
 69 >       fmt.Println("Received one login request")
 70 >       //return token (text file)
 71 >       w.Header().Set("Content-Type", "text/plain")
 72 >       w.Header().Set("Access-Control-Allow-Origin", "*")
 73 
 74 >       if r.Method == "OPTIONS" {
 75 >       >       return
 76 >       }
 77 
 78 >       decoder := json.NewDecoder(r.Body)
 79 >       var user User
 80 >       if err := decoder.Decode(&user); err != nil {
 81 >       >       http.Error(w, "Failed to read user data", http.StatusBadRequest)
 82 >       >       fmt.Printf("Failed to read user data %v\n", err)
 83 >       >       return
 84 >       }
 85 
 86 >       exists, err := checkUser(user.Username, user.Password)
 87 >       if err != nil {
 88 >       >       http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
 89 >       >       fmt.Printf("Failed to read user from Elasticsearch")
 90 >       >       return
 91 >       }
 92 
 93 >       if !exists {
 94 >       >       http.Error(w, "User does not exist", http.StatusUnauthorized)
 95 >       >       fmt.Printf("User does not exist")
 96 >       >       return
 97 >       }
 98 
 99 >       //create token
100 >       token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
101 >       >       "username": user.Username,
102 >       >       "exp":      time.Now().Add(time.Hour * 24).Unix(),
103 >       })
104 
105 >       //encrypt
106 >       tokenString, err := token.SignedString(mySigningKey)
107 >       if err != nil {
108 >       >       http.Error(w, "Failed to generate token", http.StatusInternalServerError)
109 >       >       fmt.Printf("Failed to generate token %v\n", err)
110 >       >       return
111 >       }
112 
113 >       w.Write([]byte(tokenString))
114 }
115 
116 func handlerSignup(w http.ResponseWriter, r *http.Request) {
117 >       fmt.Println("Received one signup request")
118 >       w.Header().Set("Content-Type", "text/plain")
119 >       w.Header().Set("Access-Control-Allow-Origin", "*")
120 
121 >       if r.Method == "OPTIONS" {
122 >       >       return
123 >       }
124 
125 >       decoder := json.NewDecoder(r.Body)
126 >       var user User
127 >       if err := decoder.Decode(&user); err != nil {
128 >       >       http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
129 >       >       fmt.Printf("Cannot decode user data from client %v\n", err)
130 >       >       return
131 >       }
132 
133 >       if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Userna    me) {
134 >       >       http.Error(w, "Invalid username or password", http.StatusBadRequest)
135 >       >       fmt.Printf("Invalid username or password\n")
136 >       >       return
137 >       }
138 
139 >       success, err := addUser(&user)
140 >       if err != nil {
141 >       >       http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
142 >       >       fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
143 >       >       return
144 >       }
145 
146 >       if !success {
147 >       >       http.Error(w, "User already exists", http.StatusBadRequest)
148 >       >       fmt.Println("User already exists")
149 >       >       return
150 >       }
151 
152 >       fmt.Printf("User added successfully: %s.\n", user.Username)
153 
154 }