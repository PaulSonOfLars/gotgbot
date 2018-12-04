package handlers

type baseHandler struct {
	Name string
}

func (h baseHandler) GetName() string {
	return h.Name
}
