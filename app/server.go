package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
  "context"
  "os"

  //mongo drivers
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  "github.com/nelkinda/health-go"
  "github.com/nelkinda/health-go/checks/mongodb"
)

var collection *mongo.Collection

func appHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Path)
    p := "." + r.URL.Path
    if p == "./" {
        p = "./public/index.html"
    }
  http.ServeFile(w, r, p)
  //fmt.Println(time.Now(), "Hello from my new fresh server")

}



func main() {
  //log.Println(os.Getenv("MONGODB_URI"))
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  client, db_err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
  if db_err != nil {
      log.Fatal(db_err)
  }
  
  db_err = client.Connect(ctx)
  if db_err != nil {
      log.Fatal(db_err)
  }
  defer cancel()
    
  databases, db_q_err := client.ListDatabaseNames(ctx, bson.M{})
  if db_err != nil {
       log.Fatal(db_q_err)
  }
  fmt.Println(databases)
  collection = client.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("COLLECTION"))
  //var KZN bson.M
  //if db_q_err = collection.FindOne(ctx, bson.M{}).Decode(&KZN); db_q_err != nil {
  //    log.Fatal(db_q_err)
  //}
  //fmt.Println(KZN)
  h := health.New(health.Health{Version: "1", ReleaseID: "1.0.0-SNAPSHOT"},mongodb.Health("MONGODB_URI", client, time.Duration(10)*time.Second, time.Duration(40)*time.Microsecond)) 

    // 2. Add the handler to your mux/server.
  http.HandleFunc("/health", h.Handler)

  
  s := &http.Server{
    Addr: ":443",
    Handler: nil, // use `http.DefaultServeMux`
  }

  http.HandleFunc("/", appHandler)

  log.Println("Started, serving on port 443")
  err := s.ListenAndServeTLS("public.crt", "private.key")

  //log.Println("Started, serving on port 8080")
  //err := s.ListenAndServe()

  if err != nil {
    log.Fatal(err.Error())
  }
}

