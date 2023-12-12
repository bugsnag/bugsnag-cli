package utils

func ThisOrThat(this, that interface{}) interface{} {
	switch this.(type) {
	case string:
		if this != "" {
			return this
		}

	default:
		if this != nil {
			return this
		}
	}

	return that
}

func ThisOrThatBool(this, that *bool) *bool {
	if this != nil {
		return this
	}

	return that
}
