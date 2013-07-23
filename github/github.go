package github

import (
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/plouc/go-github-client"
)

type GithubHandler struct {
    gh *gogithub.Github
}

func optionsRequestHandler(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Max-Age", "1000")
    w.Header().Set("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept, x-requested-with")

    return r.Method == "OPTIONS"
}

func New(gh *gogithub.Github) *GithubHandler {
    return &GithubHandler{gh}
}

func (ghh *GithubHandler) ReposHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    repos, err := ghh.gh.UserRepos("plouc")
    if err != nil { fmt.Println(err) }

    json, err := json.MarshalIndent(repos, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (ghh *GithubHandler) UsageHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    json, err := json.MarshalIndent(ghh.gh.RateLimit, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (ghh *GithubHandler) RepoHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)
    repo, err := ghh.gh.Repo(params["owner"], params["name"])
    if err != nil { fmt.Println("error:", err) }

    json, err := json.MarshalIndent(repo, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}
