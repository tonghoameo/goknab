package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// renewAccessTokenRequest
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_Token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	session, err := server.store.GetSession(ctx, uuid.MustParse(refreshPayload.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blok session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return

	}
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorret session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if session.FreshToken != req.RefreshToken {
		err := fmt.Errorf("mistmatch session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expried session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuraton)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.RegisteredClaims.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, rsp)

}
