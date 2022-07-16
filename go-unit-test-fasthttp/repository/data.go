package repository

type Data struct {
	RequestCardTopupCompensateLogID int
	CardNumber                      string
	CreditCode                      string
	CreditAmount                    string
	Comment                         string
	RequestStatus                   string
	DateCreate                      string
	DateModify                      string
	DateReview                      string
	CommentReview                   string
	RequestNumber                   string
	AdminCreateName                 string
	AdminModifyName                 string
	AdminReviewName                 string
}

type DataReposity interface {
	GetAll() ([]Data, error)
	GetById(id int) (*Data, error)
}
