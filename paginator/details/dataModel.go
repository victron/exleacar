package details

import (
	"errors"

	log "github.com/victron/simpleLogger"
)

type Data struct {
	Vin           string
	Photos        []string
	Specification [][]string
	DamageRaw     [][]LinkDescription
	Damage        []LinkDescription
	SupplierInfo  [][]string
}

type LinkDescription struct {
	Name     string
	Link     string
	Original string
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

func (data *Data) FilterDamage() error {
	damage := make([]LinkDescription, 0)
	var err error
	for n, row := range (*data).DamageRaw {
		if len(row) != 2 {
			log.Warning.Println("wrong len for row=", n)
			err = errors.New("wrong len in row")
			continue
		}
		if row[0].Link == row[1].Link {
			damage = append(damage, row[0])
		}
	}
	(*data).Damage = damage
	return err
}
