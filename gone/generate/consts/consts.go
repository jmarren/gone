package consts

const Converter = `func %s(data interface{}) (%s, error) {
	res, ok := data.(%s)
	if !ok {
		return nil, fmt.Errorf("data is not of type %s")
	}
	return res, nil
}`
