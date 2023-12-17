package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/token"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	fmt.Println("creatae account")
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println("CreateAccountParams  for user: ", authPayload.Username)
	// req.Owner depend on user.Username
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	//fmt.Println("init balance : ", arg.Balance)
	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		errCode := db.ErrorDbHandle(err)
		if errCode == db.UniqueViolation || errCode == db.ForeignKeyViolation {
			log.Println(errCode)
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// params request getAccountRequest
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

	}
	acc, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))

			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if acc.Owner != authPayload.Username {
		err := errors.New("account doesnt belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)

}

// params request listAccounts
type listAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageId - 1),
	}
	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))

			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)

}
