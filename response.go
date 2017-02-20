package main

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

type Status struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func NewResponse(code int, err error, data interface{}) Response {
	var status Status
	if err != nil {
		status = Status{code, err.Error()}
	} else {
		status = Status{code, ""}
	}

	return Response{
		Status: status,
		Data:   data,
	}
}
