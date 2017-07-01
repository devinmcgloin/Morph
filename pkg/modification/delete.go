package modification

//func DeleteImage(requestFrom model.Ref, image model.Ref) rsp.Response {
//	// checking if the user has permission to delete the item
//	valid, err := sql.Permissions(requestFrom.Id, model.CanDelete, image.Id)
//	if err != nil {
//		return rsp.Response{Code: http.StatusInternalServerError, Message: "Unable to retrieve user permissions."}
//	}
//	if !valid {
//		return rsp.Response{Code: http.StatusForbidden, Message: "User does not have permission to delete item."}
//	}
//
//	err = sql.DeleteImage(image.Id)
//	if err != nil {
//		return rsp.Response{Code: http.StatusInternalServerError,
//			Message: "Unable to delete user."}
//	}
//	return rsp.Response{Code: http.StatusAccepted}
//
//}
//
//func DeleteUser(requestFrom model.Ref, user model.Ref) rsp.Response {
//	// checking if the user has permission to delete the item
//	valid, err := sql.Permissions(requestFrom.Id, model.CanDelete, user.Id)
//	if err != nil {
//		return rsp.Response{Code: http.StatusInternalServerError,
//			Message: "Unable to retrieve user permissions."}
//	}
//	if !valid {
//		return rsp.Response{Code: http.StatusForbidden,
//			Message: "User does not have permission to delete item."}
//	}
//
//	err = sql.DeleteUser(user.Id)
//	if err != nil {
//		return rsp.Response{Code: http.StatusInternalServerError,
//			Message: "Unable to delete user."}
//	}
//	return rsp.Response{Code: http.StatusAccepted}
//
//}
//
//func deleteUser(id int64) error {
//	var image_ids []int64
//
//	err := db.Select(&image_ids, "SELECT id FROM content.images WHERE owner = ?", id)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	tx, err := db.Begin()
//	if err != nil {
//		log.Print(err)
//		return err
//	}
//
//	tx.Exec("DELETE FROM content.users WHERE id = ?", id)
//	for _, img := range image_ids {
//		tx.Exec("DELETE FROM content.images WHERE id = ?", img)
//		tx.Exec("DELETE FROM content.image_metadata WHERE image_id = ?", img)
//		tx.Exec("DELETE FROM content.image_label_bridge WHERE image_id = ?", img)
//		tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = ?", img)
//		tx.Exec("DELETE FROM content.user_favorites WHERE image_id = ?", img)
//	}
//	return tx.Commit()
//}
//
//// DeleteImage removes all keys for the given image, as well as removing it from
//// the owner. In the future it will also handle favorites and collections.
//func deleteImage(id int64) error {
//	tx, err := db.Begin()
//	if err != nil {
//		log.Print(err)
//		return err
//	}
//
//	tx.Exec("DELETE FROM content.images WHERE id = ?", id)
//	tx.Exec("DELETE FROM content.image_metadata WHERE image_id = ?", id)
//	tx.Exec("DELETE FROM content.image_label_bridge WHERE image_id = ?", id)
//	tx.Exec("DELETE FROM content.image_tag_bridge WHERE image_id = ?", id)
//	tx.Exec("DELETE FROM content.user_favorites WHERE image_id = ?", id)
//	return tx.Commit()
//}
