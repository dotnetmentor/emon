package main

type apiResponse struct {
	Status string               `json:"status"`
	Checks map[string]apiChecks `json:"checks"`
}

type apiChecks map[string]*check
