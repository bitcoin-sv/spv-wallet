func (s *APIUsers) CreateAddress(c *gin.Context) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	// Parse request body
	var req api.ModelsPaymailAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	fullAddress := req.Alias + "@" + req.Domain

	newAddress := &addressesmodels.NewAddress{
		UserID:             userID,
		Address:            fullAddress,
		CustomInstructions: nil,
	}

	err = reqctx.Engine(c).AddressesService().Create(c.Request.Context(), newAddress)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	// Return the created Paymail address
	c.JSON(http.StatusOK, &api.ModelsPaymailAddress{
		Alias:      req.Alias,
		Domain:     req.Domain,
		Paymail:    fullAddress,
		PublicName: req.PublicName,
		AvatarURL:  req.AvatarURL,
	})
}
