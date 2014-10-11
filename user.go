package nessusgo

type User struct {
	Admin string `json:"admin"` // Temporarily a string, due to uppercase "TRUE" Json parsing issues
	Name  string `json:"name"`
}
