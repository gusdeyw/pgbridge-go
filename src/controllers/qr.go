package controllers

// func CreateQrP(c *gin.Context) {

// 	var QueryParam global_var.PMS_Request
// 	err := c.BindJSON(&QueryParam)
// 	if err != nil {
// 		helper.SendResponse(global_var.ResponseCode.InvalidDataFormat, nil, nil, c)
// 		return
// 	}

// 	if QueryParam.PgCode == "XNDT" {
// 		result, err := CreateQrXendit(QueryParam.Detail, QueryParam.PgUsername, QueryParam.PgPassword)
// 		if err != nil {
// 			helper.SendResponse(global_var.ResponseCode.InternalServerError, "", nil, c)
// 			return
// 		} else {
// 			helper.SendResponse(global_var.ResponseCode.Successfully, "", result, c)
// 			return
// 		}
// 	} else {
// 		return
// 	}
// }

// func CreateQrXendit(QueryParam global_var.PMS_Order, PgUsername, PgPassword string) (*global_var.XDNT_ResultBody, error) {

// 	var err error
// 	var RequestBodyDetail []global_var.XDNT_ItemDetail_RequestBody

// 	for _, detail := range QueryParam.ItemDetail {
// 		RequestBodyDetail = append(RequestBodyDetail, global_var.XDNT_ItemDetail_RequestBody{
// 			ReferenceID: detail.ItemCode,
// 			Name:        detail.Name,
// 			Currency:    detail.Currency,
// 			Price:       detail.Price,
// 			Quantity:    detail.Quantity,
// 			Description: detail.Description,
// 		})
// 	}

// 	RequestBody := global_var.XDNT_RequestBody{
// 		ReferenceID: QueryParam.OrderCode,
// 		Type:        "DYNAMIC",
// 		Currency:    QueryParam.Currency,
// 		Amount:      QueryParam.Amount,
// 		ExpiresAt:   QueryParam.ExpiresAt,
// 		Basket:      RequestBodyDetail,
// 	}

// 	Header := map[string]string{
// 		"api-version": "2022-07-31",
// 	}

// 	res, _ := helper.SendRequest(global_var.RequestMethod.Post, "https://api.xendit.co/qr_codes", PgUsername, PgPassword, Header, RequestBody)
// 	var Result global_var.XDNT_ResultBody

// 	err = json.Unmarshal(res, &Result)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		return &Result, nil
// 	}
// }
