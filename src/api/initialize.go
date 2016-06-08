package api

import "github.com/devinmcgloin/morph/src/api/SQL"

func InitializeDataStores() {
	SQL.SetDB()
}
