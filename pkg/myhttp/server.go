package myhttp

import (
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) StartServer(addr *string, errorLog, infoLog *log.Logger, handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  handler,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := s.httpServer.ListenAndServe()
	errorLog.Fatal(err)
}
