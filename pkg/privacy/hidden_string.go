package privacy

type HiddenString string

func (s HiddenString) String() string {
	return "***"
}

func (s HiddenString) MarshalJSON() ([]byte, error) {
	return []byte(`"***"`), nil
}

func (s HiddenString) MarshalText() ([]byte, error) {
	return []byte("***"), nil
}
