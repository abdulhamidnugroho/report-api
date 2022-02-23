package controllers

import (
	"net/http"

	"github.com/abdulhamidnugroho/go-full/api/models"
	"github.com/abdulhamidnugroho/go-full/api/responses"
)

func (server *Server) GetMerchant(w http.ResponseWriter, r *http.Request) {
	merchant := models.Merchant{}

	merchants, err := merchant.GetAllMerchant(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, merchants)
}

func (server *Server) GetTransactionReport(w http.ResponseWriter, r *http.Request) {
	transaction := models.Transaction{}

	reports, err := transaction.GetAllTransaction(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, reports)

}
