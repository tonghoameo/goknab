package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/binbomb/goapp/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeaderKey) == 0 {
			err := errors.New("authorization header is not provider")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//	Host: example.com
		//	accept-encoding: gzip, deflate
		//	Accept-Language: en-us
		//	fOO: Bar
		//	foo: two
		//
		// then
		//
		//	Header = map[string][]string{
		//		"Accept-Encoding": {"gzip, deflate"},
		//		"Accept-Language": {"en-us"},
		//		"Foo": {"Bar", "two"},
		//	}
		/*
			eaxmples in here [Authorization: Bearer v2.local.BiS9qrXCgAfuOkAJ0DPabsRnBA17DMIqivMdOjtp0pl7-JFtkoyNmEPJXvMOhvivH37XgNNL92ZsIhp0qY8HKYK6PdcUpGaOMdXue7Wl7N2eFlqVdbidMDaTHrwPOMXyleMJqk3_V5t6HyXpEHK1JlL8VhwvDqq6OYC_uySHBxKUliYy6JuvTelL-pIQ.bnVsbA]
			map[Authorization][0] = Bearer
			map[Authorization][1] = token= v2.local.BiS9qrXCgAfuOkAJ0DPabsRnBA17DMIqivMdOjtp0pl7-JFtkoyNmEPJXvMOhvivH37XgNNL92ZsIhp0qY8HKYK6PdcUpGaOMdXue7Wl7N2eFlqVdbidMDaTHrwPOMXyleMJqk3_V5t6HyXpEHK1JlL8VhwvDqq6OYC_uySHBxKUliYy6JuvTelL-pIQ.bnVsbA
		*/
		fields := strings.Fields(authorizationHeader)
		fmt.Printf(" len fileds %d  values fileds \n", len(fields), fields)
		if len(fields) < 2 {
			// error in here
			err := errors.New("invalid authorization header	format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return

		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
