package permissions

type Permission string

const (
	CanEdit   = Permission("can_edit")
	CanDelete = Permission("can_delete")
	CanView   = Permission("can_view")
)
