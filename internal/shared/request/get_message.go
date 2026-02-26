package request

type PaginationInput struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=50"`
}