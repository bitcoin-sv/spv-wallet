package arc

import (
	"net/http"
	"net/url"
	"strings"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
)

const callbackPath = "/arc/broadcast/callback"
const bearerSchema = "Bearer "

// RegisterCallback registers the ARC callback handler and sets up final URL sent to ARC during broadcast.
func (s *Service) RegisterCallback(handler chainmodels.TXInfoHandler, router *gin.Engine) {
	if s.arcCfg.Callback == nil {
		s.logger.Info().Msg("Skipping ARC callback registration as it is not configured")
		return
	}
	hostURL, err := url.Parse(s.arcCfg.Callback.URL)
	if err != nil {
		panic(spverrors.Wrapf(err, "failed to parse ARC callback URL: %s", s.arcCfg.Callback.URL))
	}

	hostURL.Path = callbackPath
	s.arcCfg.Callback.URL = hostURL.String()

	router.POST(callbackPath, func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrMissingAuthHeader, &s.logger)
		}
		if !strings.HasPrefix(authHeader, bearerSchema) {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInvalidOrMissingToken, &s.logger)
		}

		providedToken := authHeader[len(bearerSchema):]
		if providedToken != s.arcCfg.Callback.Token {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInvalidToken, &s.logger)
		}

		var callbackResp chainmodels.TXInfo
		err = c.Bind(&callbackResp)
		if err != nil {
			spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, &s.logger)
			return
		}

		err = handler.Handle(c.Request.Context(), callbackResp)
		if err != nil {
			s.logger.Err(err).Any("TxInfo", callbackResp).Msgf("failed to update transaction in ARC broadcast callback handler")
		}

		c.Status(http.StatusOK)
	})
}
