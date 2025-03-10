package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppsCreate(w, r, s.store)
	}
}

func (s *Server) getApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppsGet(w, r, s.store)
	}
}

func (s *Server) deleteApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppDelete(w, r, s.store)
	}
}

func (s *Server) putApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppEdit(w, r, s.store)
	}
}

func (s *Server) postAppDeploy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppDeploy(w, r)
	}
}

func (s *Server) putAppStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppStop(w, r)
	}
}

func (s *Server) putAppStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppStart(w, r)
	}
}

func (s *Server) getAppLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppLogs(w, r)
	}
}

func (s *Server) getAppFailedLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppFailedLogs(w, r)
	}
}

func (s *Server) putAppRestart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppRestart(w, r)
	}
}

func (s *Server) putAppRebuild() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppRebuild(w, r)
	}
}

func (s *Server) putAppRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppRun(w, r)
	}
}
