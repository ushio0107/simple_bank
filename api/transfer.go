/*
These are the handlers of the transfer API, which are the interface between the client and the database.
*/
package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	// Check the content of json is as required of binding.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check both FromAccountID and ToAccountID are valid.
	if !s.validAccount(ctx, req.FromAccountID, req.Currency) || !s.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	res, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	ac, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if ac.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", ac.ID, ac.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
