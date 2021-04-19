package controllers

import "github.com/funthere/komodo/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")
	rs := s.Router.PathPrefix("/sellers").Subrouter()

	// seller
	rs.HandleFunc("/products", middlewares.SetMiddlewareJSON(s.GetProducts)).Methods("GET")

}
