package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"

	"github.com/satori/uuid"

	"github.com/rjeczalik/instrumenter/cmd/internal/intercept"
)

func die(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

var addr = "127.0.0.1:8484"

func main() {
	itc := &interceptor{
		errors: make(map[string]*intercept.Error),
	}

	http.HandleFunc("/new", itc.New)
	http.HandleFunc("/", itc.Show)

	log.Println("interceptor listening on", addr, "...")

	if err := http.ListenAndServe(addr, nil); err != nil {
		die(err)
	}
}

type interceptor struct {
	mu     sync.Mutex
	errors map[string]*intercept.Error
}

func (itc *interceptor) New(w http.ResponseWriter, r *http.Request) {
	var req intercept.NewRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Error == nil {
		http.Error(w, "invalid or empty error payload", http.StatusBadRequest)
	}

	req.Error.ID = uuid.NewV4().String()

	itc.mu.Lock()
	itc.errors[req.Error.ID] = req.Error
	itc.mu.Unlock()

	json.NewEncoder(w).Encode(&intercept.NewResponse{ID: req.Error.ID})
}

func (itc *interceptor) Show(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var errors []*intercept.Error
	var stacktrace bool
	var path bool

	if b, err := strconv.ParseBool(query.Get("stacktrace")); err == nil {
		stacktrace = b
	}

	if b, err := strconv.ParseBool(query.Get("path")); err == nil {
		stacktrace = b
	}

	if id := query.Get("id"); id != "" {
		itc.mu.Lock()
		err, ok := itc.errors[id]
		itc.mu.Unlock()

		if !ok {
			http.Error(w, fmt.Sprintf("error with %q id was not found"), http.StatusNotFound)
			return
		}

		errors, stacktrace, path = append(errors, err), true, true
	} else {
		errors = itc.sortedErrors()
	}

	render(w, errors, stacktrace, path)
}

func (itc *interceptor) sortedErrors() []*intercept.Error {
	itc.mu.Lock()
	errors := make([]*intercept.Error, 0, len(itc.errors))
	for _, err := range itc.errors {
		errors = append(errors, err)
	}
	itc.mu.Unlock()
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].CreatedAt.Before(errors[j].CreatedAt)
	})
	return errors
}

func render(w http.ResponseWriter, errors []*intercept.Error, stacktrace, path bool) {
	fmt.Fprintln(w, "<code>")
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	fmt.Fprint(tw, "ID\tTIME\tUSERNAME\tTYPE\tMESSAGE")

	if stacktrace {
		fmt.Fprint(tw, "\tSTACKTRACE")
	}
	if path {
		fmt.Fprint(tw, "\tPATH")
	}

	fmt.Fprintln(tw)

	for _, err := range errors {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s", err.ID, err.CreatedAt,
			err.Username, err.Type, err.Message)

		if stacktrace {
			fmt.Fprintf(tw, "\t%s", err.Stacktrace)
		}
		if path {
			fmt.Fprintf(tw, "\t%s", err.Path)
		}

		fmt.Fprintln(tw)
	}

	tw.Flush()

	fmt.Fprintln(w, "</code>")
}
