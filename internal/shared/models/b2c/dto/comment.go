package dto

type CommentRequest struct {
	Text string `json:"text" binding:"required,min=10,max=1000"`
}

func (req *CommentRequest) Validate() error {
	if len(req.Text) < 10 || len(req.Text) > 1000 {
		return ErrBadRequest
	}
	return nil
}

type CommentResponse struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	Author Author `json:"author"`
}

type Author struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	AvatarURL string `json:"avatar_url"`
}
