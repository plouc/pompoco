package projects

import (
	"fmt"
	"net/http"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type ProjectManager struct {
	collection *mgo.Collection
	decoder    *schema.Decoder
}

type Project struct {
    Id          bson.ObjectId `bson:"_id,omitempty"json:"_id,omitempty"schema:"_id"`
    Name        string        `json:"name"schema:"name"`
    Description string        `json:"description"schema:"description"`
}

func (p *Project) String() string {
    b, _ := json.Marshal(p)
    return string(b)
}

func NewProjectManager(mgo_session *mgo.Session) *ProjectManager {
	return &ProjectManager{
		mgo_session.DB("pompoco").C("projects"),
		schema.NewDecoder(),
	}
}

func optionsRequestHandler(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Max-Age", "1000")
    w.Header().Set("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept, x-requested-with")

    return r.Method == "OPTIONS"
}

func (pm *ProjectManager) getAndWriteProjects(w http.ResponseWriter) {
    projects := []*Project{}
    err := pm.collection.Find(nil).All(&projects)
    if err != nil {
        panic(err)
    }

    json, err := json.MarshalIndent(projects, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}

func (pm *ProjectManager) ProjectHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)

    project := new(Project)

    err := pm.collection.FindId(bson.ObjectIdHex(params["id"])).One(&project)
    if err != nil {
        fmt.Println(err)
    }

    json, err := json.MarshalIndent(project, "", "  ")
    if err != nil {
        fmt.Println("error:", err)
    }

    fmt.Fprint(w, string(json))
}

func (pm *ProjectManager) ProjectsHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }
    
    pm.getAndWriteProjects(w)
}

func (pm *ProjectManager) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    err := r.ParseForm()
    if err != nil { fmt.Println("error:", err) }

    project := new(Project)
    err = pm.decoder.Decode(project, r.PostForm)
    if err != nil {
        panic(err)
    }

    err = pm.collection.Insert(&project)
    if err != nil {
        panic(err)
    }

    pm.getAndWriteProjects(w)
}

func (pm *ProjectManager) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }  

    err := r.ParseForm()
    if err != nil { fmt.Println("error:", err) }

    params := mux.Vars(r)

    project := new(Project)
    err = pm.decoder.Decode(project, r.PostForm)

    err = pm.collection.UpdateId(bson.ObjectIdHex(params["id"]), project)
    if err != nil {
        panic(err)
    }

    pm.getAndWriteProjects(w)
}

func (pm *ProjectManager) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
    if optionsRequestHandler(w, r) { return }

    params := mux.Vars(r)
    err := pm.collection.RemoveId(bson.ObjectIdHex(params["id"]))
    if err != nil { fmt.Println("error:", err) }

    pm.getAndWriteProjects(w)
}
