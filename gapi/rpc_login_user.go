package gapi

import (
	"context"
	"database/sql"

	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/pb"
	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/binbomb/goapp/simplebank/valid"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	violations := validateLoginUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, " find not found user")
	}

	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wrong password")
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.GetUsername(), server.config.AccessTokenDuraton)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to create accessed token ")

	}
	// add more refreshToken
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuraton,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}
	metadata := server.extractMetada(ctx)
	arg := db.CreateSessionParams{
		ID:         uuid.MustParse(refreshPayload.RegisteredClaims.ID),
		Username:   user.Username,
		FreshToken: refreshToken,
		UserAgent:  metadata.UserAgent,
		ClientIp:   metadata.ClientIp,
		IsBlocked:  false,
		ExpiresAt:  refreshPayload.RegisteredClaims.ExpiresAt.Time,
	}

	session, err := server.store.CreateSession(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}
	rsp := &pb.LoginUserResponse{
		User:                   convertUser(user),
		SessionId:              session.ID.String(),
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		Access_TokenExpiresAt:  timestamppb.New(accessPayload.RegisteredClaims.ExpiresAt.Time),
		Refresh_TokenExpiresAt: timestamppb.New(refreshPayload.RegisteredClaims.ExpiresAt.Time),
	}

	return rsp, nil
}
func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := valid.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := valid.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
