package main

import (
    "encoding/json"
    "net/http"
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
    "github.com/emicklei/go-restful"
    "github.com/myk4040okothogodo/Gorestful/dbutils"
  )


// DB Driver visible to whole program
var DB *sql.DB
// TrainResource is the model for holding rail information
type UserResource struct {
    ID int
    UserName string
    EmailAddress string
    Age int
    MaritalStatus bool
}


// Register adds paths and routes to container
func (u *UserResource) Register(container *restful.Container) {
    us := new(restful.WebService)
    us.Path("/v1/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON) // you can specify this per route as

    us.Route(us.GET("/{user-id}").To(u.getUser))
    us.Route(us.POST("").To(u.createUser))
    us.Route(us.DELETE("/{user-id}").To(u.removeUser))
    container.Add(us)
}


// GET http://localhost:8000/v1/users/1
func (u UserResource) getUser(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("user-id")
    err := DB.QueryRow("select ID, USER_NAME,EMAIL_ADDRESS, AGE, MARITAL_STATUS FROM user where id=?", id).Scan(&u.ID, &u.UserName,&u.EmailAddress,&u.Age, &u.MaritalStatus)
    if err != nil {
        log.Println(err)
    response.AddHeader("Content-Type", "text/plain")
    response.WriteErrorString(http.StatusNotFound, "User could not be found.")
    } else {
        response.WriteEntity(u)
    }
}

// POST http://localhost:8000/v1/users
func (u UserResource) createUser(request *restful.Request, response *restful.Response) {
    log.Println(request.Request.Body)
    decoder := json.NewDecoder(request.Request.Body)
    var us UserResource
    err := decoder.Decode(&us)
    log.Println(us.UserName,us.EmailAddress,us.Age, us.MaritalStatus)
    // Error handling is obvious here. So omitting...
    statement, _ := DB.Prepare("insert into user (USER_NAME, EMAIL_ADDRESS, AGE, MARITAL_STATUS) values (?, ?, ?, ?)")
    result, err := statement.Exec(us.UserName,us.EmailAddress,us.Age, us.MaritalStatus)
    if err == nil {
        newID, _ := result.LastInsertId()
        us.ID = int(newID)
        response.WriteHeaderAndEntity(http.StatusCreated, us)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusInternalServerError,err.Error())
    }
}

// DELETE http://localhost:8000/v1/users/1
func (u UserResource) removeUser(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("user-id")
    statement, _ := DB.Prepare("delete from user where id=?")
    _, err := statement.Exec(id)
    if err == nil {
        response.WriteHeader(http.StatusOK)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusInternalServerError,err.Error())
    }
}



func main() {

    // Connect to Database
    db, err := sql.Open("sqlite3", "./userapi.db")
    if err != nil {
        log.Println("Driver creation failed!")
    }

    // Create tables
    dbutils.Initialize(db)
    wsContainer := restful.NewContainer()
    wsContainer.Router(restful.CurlyRouter{})
    u := UserResource{}
    u.Register(wsContainer)
    log.Printf("Spinning up Server on localhost:8000")
    server := &http.Server{Addr: ":8000", Handler: wsContainer}
    log.Fatal(server.ListenAndServe())
   
}
