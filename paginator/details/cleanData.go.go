package details

import (
	"errors"

	log "github.com/victron/simpleLogger"
)

func FilterDamage(data [][]LinkDescription) ([]LinkDescription, error) {
	damage := make([]LinkDescription, 0)
	var err error
	for n, row := range data {
		if len(row) != 2 {
			log.Warning.Println("wrong len for row=", n)
			err = errors.New("wrong len in row")
			continue
		}
		if row[0].Link == row[1].Link {
			damage = append(damage, row[0])
		}
	}
	return damage, err
}
