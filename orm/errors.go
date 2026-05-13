package orm

type ormErr struct {
	inner error
}

func (e ormErr) Error() string {
	return e.inner.Error()
}

func (e ormErr) Unwrap() error {
	return e.inner
}

func WrapErr(err error) error {
	return ormErr{inner: err}
}

func IsBoilErr(err error) bool {
	_, ok := err.(ormErr)
	return ok
}
