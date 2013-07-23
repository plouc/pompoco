package main

import (
	"fmt"
    "io/ioutil"
    "net/http"

	"github.com/plouc/go-gitlab-client"
    "github.com/plouc/go-github-client"
    "github.com/plouc/go-jira-client"
    "github.com/plouc/pompoco/events"
    "github.com/plouc/pompoco/projects"
    "github.com/plouc/pompoco/users"
    "github.com/plouc/pompoco/gitlab"
    "github.com/plouc/pompoco/github"
	
    "github.com/op/go-logging"
	"encoding/json"
	"github.com/gorilla/mux"
    "github.com/gorilla/schema"
    "labix.org/v2/mgo"
    es_api "github.com/mattbaird/elastigo/api"
    "launchpad.net/goyaml"
)

var (
    gl          *gogitlab.Gitlab
    gh          *gogithub.Github
    jira        *gojira.Jira
    settings    *Settings
    mgo_session *mgo.Session
    users_c     *mgo.Collection
    decoder     = schema.NewDecoder()
    log         = logging.MustGetLogger("pompoco")
)

type Settings struct {
    LogFile string `yaml:"log_file"`
    Gitlab struct{
        Host    string `yaml:"host"`
        ApiPath string `yaml:"api_path"`
        Token   string `yaml:"token"`
    }
    Jira struct{
        Host         string `yaml:"host"`
        ApiPath      string `yaml:"api_path"`
        ActivityPath string `yaml:"activity_path"`
        Login        string `yaml:"login"`
        Password     string `yaml:"password"`
    }
    Elasticsearch struct{
        Host string `yaml:"host"` 
    }
    Mongo struct{
        Host string `yaml:"host"`
        Port string `yaml:"port"` 
    }
}

func optionsRequestHandler(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Max-Age", "1000")
    w.Header().Set("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept, x-requested-with")

    return r.Method == "OPTIONS"
}

func jiraTimelineHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    log.Debug("jiraTimelineHandler")

    activity, err := jira.UserActivity("raphael.benitte")
    if err != nil { fmt.Println("error:", err) }

    json, err := json.MarshalIndent(activity, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    log.Debug("settingsHandler")
    
    json, err := json.MarshalIndent(settings, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}

func projectBindingHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }  

    log.Debug("projectBindingHandler")

    params := mux.Vars(r)
    fmt.Println(params["id"])
    fmt.Println(params["object_type"])
    fmt.Println(params["object_role"])
    fmt.Println(params["object_id"])
}


func main() {

    logging.SetLevel(logging.DEBUG, "pompoco")

    // read config file
    file, err := ioutil.ReadFile("./.pompoco.yml")
    if err != nil { panic(err) }

    // parse config file
    settings := new(Settings)
    err = goyaml.Unmarshal(file, &settings)
    if err != nil { panic(err) }

    es_api.Domain = settings.Elasticsearch.Host

    // init mongo
    log.Debug("Connecting to mongodb")
    mgo_session, err := mgo.Dial(fmt.Sprintf("%s:%s", settings.Mongo.Host, settings.Mongo.Port))
    if err != nil { panic(err) }
    defer mgo_session.Close()

    // create clients
    log.Debug("Creating jira client")
    jira = gojira.NewJira(
        settings.Jira.Host,
        settings.Jira.ApiPath,
        settings.Jira.ActivityPath,
        &gojira.Auth{settings.Jira.Login, settings.Jira.Password,},
    )

    log.Debug("Creating gitlab client")
    gl = gogitlab.NewGitlab(
        settings.Gitlab.Host,
        settings.Gitlab.ApiPath,
        settings.Gitlab.Token,
    )

    log.Debug("Creating github client")
    gh = gogithub.NewGithub()

    r := mux.NewRouter()

    // settings
    r.HandleFunc("/settings", settingsHandler)

    em := events.NewEventManager(gl, jira)
    r.HandleFunc("/events",          em.EventsHandler)
    r.HandleFunc("/events/sync",     em.EventsSyncHandler)
    r.HandleFunc("/events/autocomp", em.EventsAutocompleteHandler)

    // users
    um := users.NewUserManager(mgo_session, gh, gl, jira)
    r.HandleFunc("/users",                 um.UsersHandler)
    r.HandleFunc("/user/{id}",             um.UserHandler).Methods("GET", "OPTIONS")
    r.HandleFunc("/user",                  um.CreateUserHandler).Methods("POST", "OPTIONS")
    r.HandleFunc("/user/{id}",             um.UpdateUserHandler).Methods("PUT", "OPTIONS")
    r.HandleFunc("/user/{id}",             um.DeleteUserHandler).Methods("DELETE", "OPTIONS")
    r.HandleFunc("/user/{id}/consolidate", um.ConsolidateUserHandler).Methods("GET", "OPTIONS")

    // projects
    pm := projects.NewProjectManager(mgo_session)
    r.HandleFunc("/projects",     pm.ProjectsHandler)
    r.HandleFunc("/project",      pm.CreateProjectHandler).Methods("POST", "OPTIONS")
    r.HandleFunc("/project/{id}", pm.ProjectHandler).Methods("GET", "OPTIONS")
    r.HandleFunc("/project/{id}", pm.UpdateProjectHandler).Methods("PUT", "OPTIONS")
    r.HandleFunc("/project/{id}", pm.DeleteProjectHandler).Methods("DELETE", "OPTIONS")

    r.HandleFunc("/project/{id}/bind", projectBindingHandler)

    // jira
    r.HandleFunc("/jira/timeline", jiraTimelineHandler)

    // gitlab
    glh := gitlab.New(gl)
    r.HandleFunc("/gitlab/timeline",              glh.TimelineHandler)
    r.HandleFunc("/gitlab/projects",              glh.ProjectsHandler)
    r.HandleFunc("/gitlab/project/{id}",          glh.ProjectHandler)
    r.HandleFunc("/gitlab/project/{id}/tags",     glh.ProjectTagsHandler)
    r.HandleFunc("/gitlab/project/{id}/branches", glh.ProjectBranchesHandler)

    // github
    ghh := github.New(gh)
    r.HandleFunc("/github/usage",               ghh.UsageHandler)
    r.HandleFunc("/github/repos",               ghh.ReposHandler)
    r.HandleFunc("/github/repo/{owner}/{name}", ghh.RepoHandler)

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
