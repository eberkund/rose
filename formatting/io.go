package formatting

type encoder interface {
	Encode(v interface{}) error
}

type decoder interface {
	Decode(v interface{}) error
}

func format(encoder encoder, decoder decoder) error {
	var decoded interface{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return err
	}
	return encoder.Encode(decoded)
}
