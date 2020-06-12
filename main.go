package main

import (
	"net/http"

	"github.com/kelseyhightower/run"
)

func main() {
	run.Notice("Starting badger service...")

	http.HandleFunc("/test/build/status", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			badRequest(w)
			return
		}

		project := r.FormValue("project")
		if project != "test" {
			badRequest(w)
			return
		}

		switch id {
		case "success":
			success(w)
		case "failure":
			failure(w)
		case "working":
			working(w)
		default:
			unknown(w)
		}
	})

	http.HandleFunc("/build/status", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			badRequest(w)
			return
		}

		project := r.FormValue("project")
		if project == "" {
			badRequest(w)
			return
		}

		status, err := getBuildStatus(project, id)
		if err != nil {
			run.Error(r, err)
			internalError(w)
			return
		}

		switch status {
		case "SUCCESS":
			success(w)
		case "FAILURE":
			failure(w)
		case "WORKING":
			working(w)
		default:
			unknown(w)
		}
	})

	run.Fatal(run.ListenAndServe(nil))
}

func success(w http.ResponseWriter) { ok(w, svgSuccess, svgSuccessHash) }
func failure(w http.ResponseWriter) { ok(w, svgFailure, svgFailureHash) }
func unknown(w http.ResponseWriter) { ok(w, svgUnknown, svgUnknownHash) }
func working(w http.ResponseWriter) { ok(w, svgWorking, svgWorkingHash) }

func ok(w http.ResponseWriter, svg []byte, etag string) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", etag)
	w.WriteHeader(http.StatusOK)
	w.Write(svg)
}

func badRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", svgUnknownHash)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(svgUnknown)
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", svgUnknownHash)
	w.WriteHeader(http.StatusNotFound)
	w.Write(svgUnknown)
}

func internalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", svgUnknownHash)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(svgUnknown)
}
