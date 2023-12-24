package panic

func Panic(v any) {
	panic(v)
}

func OnCond(cond bool, desc any) {
	if cond {
		panic(desc)
	}
}

func OnError(err error) {
	if err != nil {
		panic(err)
	}
}
