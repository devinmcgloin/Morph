package endpoint

import (
	"database/sql"
	"math/rand"
	"reflect"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/gorilla/schema"
)

func init() {
	nullString, nullBool, nullInt64, nullFloat64 := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}

	decoder.RegisterConverter(nullString, ConvertSQLNullString)
	decoder.RegisterConverter(nullBool, ConvertSQLNullBool)
	decoder.RegisterConverter(nullInt64, ConvertSQLNullInt64)
	decoder.RegisterConverter(nullFloat64, ConvertSQLNullFloat64)
}

var decoder = schema.NewDecoder()

func getNewIID() uint64 {
	IID := uint64(rand.Int63n(10000))
	for SQL.ExistsIID(IID) || IID == 0 {
		IID = uint64(rand.Int63n(10000))
	}
	return IID
}

func getNewSID() uint64 {
	SID := uint64(rand.Int63n(10000))
	for SQL.ExistsSID(SID) || SID == 0 {
		SID = uint64(rand.Int63n(10000))
	}
	return SID
}

func ConvertSQLNullString(value string) reflect.Value {
	v := sql.NullString{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func ConvertSQLNullBool(value string) reflect.Value {
	v := sql.NullBool{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func ConvertSQLNullInt64(value string) reflect.Value {
	v := sql.NullInt64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func ConvertSQLNullFloat64(value string) reflect.Value {
	v := sql.NullFloat64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}
