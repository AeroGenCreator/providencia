package types

type Employee struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	HireDate     string  `json:"hireDate"`
	DailyPayment float64 `json:"dailyPayment"`
	Phone        int     `json:"phone"`
	Email        string  `json:"email"`
	Nss          string  `json:"nss"`
	Curp         string  `json:"curp"`
}

type EmployeeSlice struct {
	Records     []*Employee `json:"records"`
	HasNextPage bool        `json:"hasNextPage"`
}

type OnlyId struct {
	Id int `json:"id"`
}
