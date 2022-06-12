package helpers

// This trick helps to create constant errors Eg
// const SomeError = Error("some error text")
type Error string

func (e Error) Error() string {
	return string(e)
}
