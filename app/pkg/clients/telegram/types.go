package telegram

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:" "`
}

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
