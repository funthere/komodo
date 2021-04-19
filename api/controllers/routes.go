package controllers

import "github.com/funthere/komodo/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")
	rs := s.Router.PathPrefix("/sellers").Subrouter()
	rb := s.Router.PathPrefix("/buyers").Subrouter()

	// Login Route
	rs.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// seller
	rs.HandleFunc("/products", middlewares.SetMiddlewareJSON(s.GetSellerProducts)).Methods("GET")
	rs.HandleFunc("/products", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateProduct))).Methods("POST")

	// buyers
	rb.HandleFunc("/products", middlewares.SetMiddlewareJSON(s.GetProducts)).Methods("GET")
}
