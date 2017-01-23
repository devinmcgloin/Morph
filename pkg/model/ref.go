package model

import "fmt"

type RString string

const (
	Images      RString = "images"
	Users       RString = "users"
	Collections RString = "collections"

	CollectionsIn RString = "collections_in"
	Downloads     RString = "downloads"
	Featured      RString = "featured"
	Views         RString = "views"
	Purchases     RString = "purchases"
	Owner         RString = "owner"
	CanView       RString = "can_view"
	CanEdit       RString = "can_edit"
	CanDelete     RString = "can_delete"
	Admin         RString = "admins"
	Purchased     RString = "purchased"
	Downloaded    RString = "downloaded"
	Seen          RString = "seen"
	FollowedBy    RString = "followed_by"
	Followed      RString = "followed"
	Favorited     RString = "favorited"
	FavoritedBy   RString = "favorited_by"
	Stream        RString = "stream"
	ViewType      RString = "view_type"
	New           RString = "new"
	Location      RString = "location"
)

type Ref struct {
	Collection RString
	ShortCode  string
}

func (ref Ref) GetTag() string {
	return fmt.Sprintf("%s:%s", ref.Collection, ref.ShortCode)
}

func (ref Ref) GetRString(t RString) string {
	return fmt.Sprintf("%s:%s:%s", ref.Collection, ref.ShortCode, t)
}

func (ref Ref) GetRSet(t RString) string {
	return fmt.Sprintf("%s:%s", ref.Collection, t)
}

func (ref Ref) Valid(typeOpts ...RString) bool {
	if ref.ShortCode == "" {
		return false
	}

	for _, t := range typeOpts {
		if t == ref.Collection {
			return true
		}
	}
	return false
}
