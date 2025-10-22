func DecodeJWT(ctx context.Context, c *app.RequestContext) {
	var err error
	var req numberOperation.DecodeJWTReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(numberOperation.DecodeJWTResp)

	rawJWT := string(c.Cookie("JWT"))
	fmt.Println("RawJWT:", rawJWT)
	fmt.Println("Decoding JWT...")
	res, err := JWT.ParseJWT(rawJWT)

	if err != nil {
		resp.IsValid = false
		c.JSON(consts.StatusOK, resp)
		return
	}

	fmt.Println("Decoding JWT...:", res)
	resp.IsValid = true
	payload := make(map[string]string)
	for k, v := range res {
		payload[k] = fmt.Sprintf("%v", v)
	}
	resp.PayLoad = payload



	c.JSON(consts.StatusOK, resp)
}