package header

import "strings"

const (
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
)

type CTX interface {
	GetHeader(key string) string
}

func GetProfile(c CTX) Profile {
	return getProfile(c)
}

// Profile quanxiang`s user profile information
type Profile struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	DepartmentID string `json:"department_id"`
}

func GetDepartments(c CTX) []string {
	departmentID := c.GetHeader(_departmentID)
	return strings.Split(departmentID, ",")
}

func getProfile(c CTX) Profile {
	userID := c.GetHeader(_userID)
	userName := c.GetHeader(_userName)
	departmentID := c.GetHeader(_departmentID)

	return Profile{
		UserID:       userID,
		UserName:     userName,
		DepartmentID: strings.Split(departmentID, ",")[0],
	}
}
