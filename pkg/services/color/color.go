package color

import (
	"log"

	"github.com/devinmcgloin/clr/clr"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type ColorCatagory string

const Shade = "shade"
const SpecificColor = "specific"

type FokalColorTable struct {
	db   *sqlx.DB
	Type ColorCatagory
}

func New(db *sqlx.DB) *FokalColorTable {
	return &FokalColorTable{
		db:   db,
		Type: SpecificColor,
	}

}

func NewWithType(db *sqlx.DB, kind ColorCatagory) *FokalColorTable {
	return &FokalColorTable{
		db:   db,
		Type: kind,
	}

}

func (spc FokalColorTable) Iterate() []clr.Color {
	hexColors := []struct{ Hex string }{}
	err := spc.db.Select(&hexColors, `SELECT hex FROM colors.clr WHERE type = $1`, spc.Type)

	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return []clr.Color{}
	}

	colors := make([]clr.Color, len(hexColors))
	for i, color := range hexColors {
		colors[i] = clr.Hex{Code: color.Hex}
	}

	return colors
}

func (spc FokalColorTable) Lookup(hex string) clr.ColorSpace {
	var name string
	err := spc.db.Get(&name, `SELECT name FROM colors.clr WHERE hex = $1`, hex)
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return ""
	}
	return clr.ColorSpace(name)
}

func (spc FokalColorTable) AddColor(name, hex string) error {
	_, err := spc.db.Exec("INSERT INTO colors.clr(name, hex, type) VALUES($1, $2, $3)",
		name, hex, spc.Type)
	if err != nil {
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return errors.Wrap(err, "unable to inset new color")
	}
	return nil
}

func (spc FokalColorTable) AddColors(colors map[string]string) error {

	tx, err := spc.db.Begin()
	if err != nil {
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return errors.Wrap(err, "unable to begin transaction")
	}
	stmt, err := tx.Prepare(pq.CopyInSchema("colors", "clr", "name", "hex", "type"))
	if err != nil {
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return errors.Wrap(err, "unable to begin transaction")
	}

	for hex, name := range colors {
		_, err = stmt.Exec(name, hex, spc.Type)
		if err != nil {
			log.Println(err)
			raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})

		return err
	}

	err = stmt.Close()
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return err
	}
	return nil
}

func (spc FokalColorTable) Colors() (map[string]string, error) {
	var clrs = make(map[string]string)
	var hex []struct {
		Name string
		Hex  string
	}
	err := spc.db.Select(&hex, "SELECT name, hex FROM colors.clr WHERE type = $1", spc.Type)
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "postgresql", "module": "color"})
		return clrs, err
	}

	for _, color := range hex {
		clrs[color.Hex] = color.Name
	}
	return clrs, nil
}
