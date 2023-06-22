package entity

func isValidName(name string) error {
	if name == "" {
		return ErrorName
	}
	return nil
}
