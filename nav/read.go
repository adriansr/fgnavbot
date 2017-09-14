package nav

type Reader interface {
	OnNext(interface{})
	OnComplete()
	OnError(error)
}
