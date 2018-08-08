package rxgo

type Connectable interface {
	Connect()
}

type ConnectableObservable interface {
	Connectable
	Observable
}

func Multicast(func(o Observable) Subject) ConnectableObservable {
	return nil
}
