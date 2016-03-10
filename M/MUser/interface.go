package MUser

type user struct {
	id int
	username string
	email string
	pwd []byte
	sale []byte
	master bool
}
