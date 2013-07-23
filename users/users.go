package users

import (
	"fmt"
	"net/http"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/plouc/go-gitlab-client"
    "github.com/plouc/go-github-client"
    "github.com/plouc/go-jira-client"
)

type UserManager struct {
	collection *mgo.Collection
	decoder    *schema.Decoder
	gh         *gogithub.Github
	gl         *gogitlab.Gitlab
    jira       *gojira.Jira
}

type User struct {
    Id bson.ObjectId `bson:"_id,omitempty"json:"_id"schema:"_id"`

    // settings
    Name         string `bson:"name"json:"name"schema:"name"`
    GithubActive bool   `bson:"github_active"json:"github_active"schema:"github_active"`
    GithubName   string `bson:"github_name"json:"github_name"schema:"github_name"`
    GitlabActive bool   `bson:"gitlab_active"json:"gitlab_active"schema:"gitlab_active"`
    GitlabId     string `bson:"gitlab_id"json:"gitlab_id"schema:"gitlab_id"`
    JiraActive   bool   `bson:"jira_active"json:"jira_active"schema:"jira_active"`
    JiraUsername string `bson:"jira_username"json:"jira_username"schema:"jira_username"`

    // api profiles
    GithubProfile *gogithub.PublicUser `bson:"github_profile"json:"github_profile"`
    GitlabProfile *gogitlab.User       `bson:"gitlab_profile"json:"gitlab_profile"`
    JiraProfile   *gojira.User         `bson:"jira_profile"json:"jira_profile"`
}

func NewUserManager(
	mgo_session *mgo.Session,
	gh          *gogithub.Github,
	gl          *gogitlab.Gitlab,
	jira        *gojira.Jira) *UserManager {

	return &UserManager{
		mgo_session.DB("pompoco").C("users"),
		schema.NewDecoder(),
		gh,
		gl,
		jira,
	}
}

func optionsRequestHandler(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Max-Age", "1000")
    w.Header().Set("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept, x-requested-with")

    return r.Method == "OPTIONS"
}

func (um *UserManager) GetUser(id string) (*User, error) { 
    user := new(User)
    err := um.collection.FindId(bson.ObjectIdHex(id)).One(&user)

    return user, err
}

func (um *UserManager) consolidateProfile(user_id string) *User {
    fmt.Println("consolidateProfile: " + user_id)

    user, err := um.GetUser(user_id)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%+v\n", user)

    if user.GithubActive {
        ghUser, err := um.gh.GetUser(user.GithubName) 
        if err != nil {
            fmt.Println(err)
        }
        user.GithubProfile = ghUser
    }

    if user.GitlabActive {
        glUser, err := um.gl.User(user.GitlabId)
        if err != nil {
            fmt.Println(err)
        }
        user.GitlabProfile = glUser
    }

    if user.JiraActive {
        jiraUser, err := um.jira.User(user.JiraUsername)
        if err != nil {
            fmt.Println(err)
        }
        user.JiraProfile = jiraUser
    }

    fmt.Printf("%+v\n", user)

    err = um.collection.UpdateId(user.Id, user)
    if err != nil {
        panic(err)
    }

    return user
}

func (um *UserManager) ConsolidateUserHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)

    user := um.consolidateProfile(params["id"])

    json, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        fmt.Println("error:", err)
    }

    fmt.Fprint(w, string(json))
}

func (um *UserManager) UserHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)

    user, err := um.GetUser(params["id"])
    if err != nil {
        fmt.Println(err)
    }

    json, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        fmt.Println("error:", err)
    }

    fmt.Fprint(w, string(json))
}


func (um *UserManager) getAndWriteUsers(w http.ResponseWriter) {
    users := []*User{}
    err := um.collection.Find(nil).All(&users)
    if err != nil {
        panic(err)
    }

    json, err := json.MarshalIndent(users, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}

func (um *UserManager) UsersHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }
    
    um.getAndWriteUsers(w)
}

func (um *UserManager) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)
    fmt.Println(params["id"])

    err := um.collection.RemoveId(bson.ObjectIdHex(params["id"]))
    if err != nil { fmt.Println("error:", err) }

    um.getAndWriteUsers(w)
}

func (um *UserManager) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    err := r.ParseForm()
    if err != nil { fmt.Println("error:", err) }

    user := new(User)
    err = um.decoder.Decode(user, r.PostForm)
    if err != nil {
        panic(err)
    }

    err = um.collection.Insert(&user)
    if err != nil {
        panic(err)
    }

    um.getAndWriteUsers(w)
}

func (um *UserManager) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }  

    err := r.ParseForm()
    if err != nil { fmt.Println("error:", err) }

    params := mux.Vars(r)

    user := new(User)
    err = um.decoder.Decode(user, r.PostForm)

    err = um.collection.UpdateId(bson.ObjectIdHex(params["id"]), user)
    if err != nil {
        panic(err)
    }

    um.getAndWriteUsers(w)
}
