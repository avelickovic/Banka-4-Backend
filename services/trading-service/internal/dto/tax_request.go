package dto

type ListTaxUsersRequest struct {
	UserType  string `form:"userType"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	Page      int32  `form:"page"`
	PageSize  int32  `form:"page_size"`
}
