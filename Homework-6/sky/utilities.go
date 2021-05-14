package sky

func illegalPanic(ok bool, info string) {
	if ok {
		panic(info)
	}
}
