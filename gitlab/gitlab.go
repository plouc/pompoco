package gitlab

import (
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/plouc/go-gitlab-client"
)

type GitlabHandler struct {
    gl *gogitlab.Gitlab
}

func optionsRequestHandler(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Max-Age", "1000")
    w.Header().Set("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept, x-requested-with")

    return r.Method == "OPTIONS"
}

func New(gl *gogitlab.Gitlab) *GitlabHandler {
    return &GitlabHandler{gl}
}

func (glh *GitlabHandler) ProjectsHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    projects, err := glh.gl.Projects()
    if err != nil { fmt.Println(err) }

    json, err := json.MarshalIndent(projects, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (glh *GitlabHandler) TimelineHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    events, err := glh.gl.Activity()
    if err != nil { fmt.Println(err) }

    json, err := json.MarshalIndent(events, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (glh *GitlabHandler) ProjectHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)
    project, err := glh.gl.Project(params["id"])
    if err != nil { fmt.Println(err) }

    json, err := json.MarshalIndent(project, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (glh *GitlabHandler) ProjectBranchesHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)
    branches, err := glh.gl.RepoBranches(params["id"])
    if err != nil { fmt.Println(err) }

    json, err := json.MarshalIndent(branches, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (glh *GitlabHandler) ProjectTagsHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)
    tags, err := glh.gl.RepoTags(params["id"])
    if err != nil { fmt.Println(err) }

    json, err := json.MarshalIndent(tags, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}