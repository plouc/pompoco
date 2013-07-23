package events

import (
	"fmt"
	"net/http"
	"github.com/plouc/go-gitlab-client"
	"github.com/plouc/go-jira-client"
	es_core "github.com/mattbaird/elastigo/core"
    es_indices "github.com/mattbaird/elastigo/indices"
    es_search "github.com/mattbaird/elastigo/search"
    "encoding/json"
    "time"
    "reflect"
//    "github.com/gorilla/mux"
)

type EventManager struct {
	gl   *gogitlab.Gitlab
	jira *gojira.Jira
}

type Event struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	OccurredAt  string `json:"occurred_at"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

func (e *Event) String() string {
    s, _ := json.Marshal(e)
    return string(s)
}

func NewEventManager(gl *gogitlab.Gitlab, jira *gojira.Jira) *EventManager {
	return &EventManager{
		gl:   gl,
		jira: jira,
	}
}

func optionsRequestHandler(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Max-Age", "1000")
    w.Header().Set("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept, x-requested-with")

    return r.Method == "OPTIONS"
}

func (em *EventManager) EventsSyncHandler(w http.ResponseWriter, r *http.Request) {
	if optionsRequestHandler(w, r) { return }

	em.SyncEvents()

	em.EventsHandler(w, r)
}

/*
Handle autocompletion query
{
  "query": {
    "query_string": {
      "fields": [
        "description.autocomplete"
      ],
      "query": "raph*"
    }
  }
}
*/
func (em *EventManager) EventsAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	if optionsRequestHandler(w, r) { return }

	err := r.ParseForm()
    if err != nil { fmt.Println("error:", err) }

	q, present := r.Form["q"]
	if !present {
		fmt.Println("\nNo query specified\n")
	}
	fmt.Println("\ntype\n", reflect.TypeOf(q))
	fmt.Printf("\nQUERY %+v\n", q)

	qry := map[string]interface{}{
	  "query":map[string]interface{}{
	    "query_string":map[string]interface{}{
	      "fields":[]string{"description.autocomplete"},
	      "query":"*",
	    },
	  },
	}
	
	es_core.SearchRequest(true, "pompoco", "event", qry, "")

	out, err := es_core.SearchRequest(true, "pompoco", "event", qry, "")
	if err != nil {
		fmt.Println(err)
	}

	json, err := json.MarshalIndent(out, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}

func (em *EventManager) EventsHandler(w http.ResponseWriter, r *http.Request) {
	if optionsRequestHandler(w, r) { return }
	
	events := em.GetEvents()

    json, err := json.MarshalIndent(events, "", "  ")
    if err != nil { fmt.Println("error:", err) }

    fmt.Fprint(w, string(json))
}


func (em *EventManager) SyncEvents() {
	glActivity, err := em.gl.Activity()
	if err != nil {
		fmt.Println(err)
	}

	for _, glEvent := range(glActivity.Entries) {
		fmt.Printf("%+v\n", glEvent)
		event := Event{
			Id:          glEvent.Id,
			Type:        "gitlab",
			OccurredAt:  glEvent.Updated.Format(time.RFC3339),
			Description: glEvent.Title,
			Username:    glEvent.Author.Name,
		}
		response, _ := es_core.Index(true, "pompoco", "event", "", event)
		fmt.Printf("Index OK: %v", response.Ok)
	}
	es_indices.Flush()    

	jiraActivity, err := em.jira.UserActivity("raphael.benitte")
	if err != nil {
		fmt.Println(err)
	}
	for _, jiraEvent := range(jiraActivity.Entries) {
		event := Event{
			Id:          jiraEvent.Id,
			Type:        "jira",
			OccurredAt:  jiraEvent.Updated.Format(time.RFC3339),
			Description: jiraEvent.Title,
			Username:    jiraEvent.Author.Name,
		}
		response, _ := es_core.Index(true, "pompoco", "event", "", event)
		fmt.Printf("Index OK: %v", response.Ok)
	}
	es_indices.Flush()  
}


func (em *EventManager) GetEvents() *[]*Event {

	events :=  make([]*Event, 0)

	out, err := es_search.Search("pompoco").Type("event").Sort(es_search.Sort("occurred_at").Desc()).Size("100").Result()
	if err != nil {
		fmt.Println(err)
	}

	for _, hit := range out.Hits.Hits {
		var event *Event
		err = json.Unmarshal(hit.Source, &event)
		if err != nil {
			fmt.Println("%s", err)
		}
		events = append(events, event)
	}

	return &events
}