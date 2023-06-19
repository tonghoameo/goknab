package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
type NewUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rspUser := newUserResponse(user)
	ctx.JSON(http.StatusOK, rspUser)
}

// loginUserRequest
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID       `json:"session_id"`
	AccessToken           string          `json:"access_token"`
	AccessTokenExpiresAt  time.Time       `json:"access_Token_expires_at"`
	RefreshToken          string          `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time       `json:"refresh_Token_expires_at"`
	User                  NewUserResponse `json:"user"`
}

func newUserResponse(user db.User) NewUserResponse {

	return NewUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Printf(" %s %s \n", req.Password, user.HashedPassword)
	err = utils.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuraton)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	// add more refreshToken
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuraton)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("refreshPayload ", refreshPayload)
	arg := db.CreateSessionParams{
		ID:         uuid.MustParse(refreshPayload.RegisteredClaims.ID),
		Username:   user.Username,
		FreshToken: refreshToken,
		UserAgent:  ctx.Request.UserAgent(),
		ClientIp:   ctx.ClientIP(),
		IsBlocked:  false,
		ExpiresAt:  refreshPayload.RegisteredClaims.ExpiresAt.Time,
	}
	fmt.Println("create session")
	session, err := server.store.CreateSession(ctx, arg)
	fmt.Println("check session")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.RegisteredClaims.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.RegisteredClaims.ExpiresAt.Time,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)

}
