package storage

type Credentials struct {
	UserId   int32 // ID of the user these credentials belong to
	Username string
	Passhash string
}
