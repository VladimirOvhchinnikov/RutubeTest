package usecase

type UseCaseInterface interface {
	StartCase(firstName string, lastname string, id int) error
	SetBirthday(date string, id int) error
	SetAllUser(id int) error
	SetSub(id int, idSub string) error
}
