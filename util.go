package vnprovince

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func ignore[T any](v T, _ error) T {
	return v
}
