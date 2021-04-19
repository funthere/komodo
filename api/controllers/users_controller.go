package controllers

import (
	"net/http"

	"github.com/funthere/komodo/api/models"
	"github.com/funthere/komodo/api/responses"
)

func (server *Server) GetProducts(w http.ResponseWriter, r *http.Request) {

	product := models.Product{}

	products, err := product.FindAll(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, products)
}
