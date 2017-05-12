package dnspass

type hostMap map[string]struct{}

func (m *hostMap) UnmarshalTOML(decode func(interface{}) error) error {
	list := []string{}
	if err := decode(&list); err != nil {
		return err
	}

	*m = make(hostMap)
	for _, group := range list {
		(*m)[group] = emptyStruct
	}

	return nil
}

func (m *hostMap) isBad(group ...string) bool {
	for _, g := range group {
		if _, ok := (*m)[g]; ok {
			return true
		}
	}
	return false
}
