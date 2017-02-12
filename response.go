package main

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

type Status struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
