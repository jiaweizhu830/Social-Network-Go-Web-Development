1 package main
  2 
  3 import (
  4 >       "context"
  5 >       "fmt"
  6 
  7 >       "github.com/olivere/elastic"
  8 )
  9 
 10 const (
 11 >       POST_INDEX = "post"
 12 >       USER_INDEX = "user"
 13 >       ES_URL     = "http://10.128.0.2:9200"
 14 )
 15 
 16 func main() {
 17 >       client, err := elastic.NewClient(elastic.SetURL(ES_URL))
 18 >       if err != nil {
 19 >       >       panic(err)
 20 >       }
 21 >       exists, err := client.IndexExists(POST_INDEX).Do(context.Background())
 22 >       if err != nil {
 23 >       >       panic(err)
 24 >       }
 25 >       if !exists {
 26 >       >       mapping := `{
 27                         "mappings": {
 28                                 "properties": {
 29                                         "user":     { "type": "keyword", "index": false },
 30                                         "message":  { "type": "keyword", "index": false },
 31                                         "location": { "type": "geo_point" },
 32                                         "url":      { "type": "keyword", "index": false },
 33                                         "type":     { "type": "keyword", "index": false },
 34                                         "face":     { "type": "float" }
 35                                 }
 36                         }
 37                 }`
 38 >       >       _, err := client.CreateIndex(POST_INDEX).Body(mapping).Do(context.Background())
 39 >       >       if err != nil {
 40 >       >       >       panic(err)
 41 >       >       }
 42 >       }
 43 
 44 >       exists, err = client.IndexExists(USER_INDEX).Do(context.Background())
 45 >       if err != nil {
 46 >       >       panic(err)
 47 >       }
 48 
 49 >       if !exists {
 50 >       >       mapping := `{
 51                         "mappings": {
 52                                 "properties": {
 53                                         "username": {"type": "keyword"},
 54                                         "password": {"type": "keyword", "index": false},
 55                                         "age":      {"type": "long", "index": false},
 56                                         "gender":   {"type": "keyword", "index": false}
 57                                 }
 58                         }
 59                 }`
 60 >       >       _, err = client.CreateIndex(USER_INDEX).Body(mapping).Do(context.Background())
 61 >       >       if err != nil {
 62 >       >       >       panic(err)
 63 >       >       }
 64 >       }
 65 
 66 >       fmt.Println("Post index is created.")
 67 }
