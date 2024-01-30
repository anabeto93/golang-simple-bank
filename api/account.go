package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/anabeto93/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD GHS NGN"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Balance: 0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
	return
}

type listAccountsRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize, // 0 * 5 = 0, 1 * 5 = 5, 2 * 5 = 10
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	Balance int64 `json:"balance" binding:"required,min=1"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	// first get the account id
	var idReq getAccountRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// next validate the request body
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// next get the account
	account, err := server.store.GetAccount(ctx, idReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// next update the account
	arg := db.UpdateAccountParams{
		ID: idReq.ID,
		Balance: req.Balance,
	}

	account, err = server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// finally return the updated account
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	// first get the account id
	var idReq getAccountRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// next fetch the account
	account, err := server.store.GetAccount(ctx, idReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// log the error
			fmt.Sprintf("account with id %d not found\n", idReq.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// next delete the account
	err = server.store.DeleteAccount(ctx, account.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// finally return a no-content response
	ctx.Status(http.StatusNoContent)
}
