package telegram

type RouteNotFoundError struct {
	Name string
}

func (r RouteNotFoundError) Error() string {
	return r.Name + " is not found"
}

type DataNotFound struct {
}

func (d DataNotFound) Error() string {
	return "your requesting data is not found"
}

type UnknownError struct {
}

func (d UnknownError) Error() string {
	return "Unknown error"
}
