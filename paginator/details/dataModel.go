package details

import (
	"errors"
)

type Data struct {
	Vin           string
	Condition     string
	Photos        []string
	Specification [][]string
	Comments      []string
	Damage        []LinkDescription
	SupplierInfo  [][]string
	RawData       struct {
		Damage [][]string
	}
}

type LinkDescription struct {
	Name string
	Link string
}

// get vin
func (data *Data) GetVin() (string, error) {
	for _, row := range (*data).Specification {
		for _, col := range row {
			if col == "VIN number" {
				if len(row) != 2 {
					return "", errors.New("wrong row len with VIN")
				}
				return row[1], nil
			}
		}
	}
	return "", errors.New("VIN not found")
}
