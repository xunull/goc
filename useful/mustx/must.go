package mustx

func StringNotNull(str string) {
	if str == "" {
		panic(ErrorStringMustNotNull)
	}
}
