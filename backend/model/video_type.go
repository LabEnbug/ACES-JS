package model

type VideoType struct {
	// Id, TypeName
	Id       uint   `default:"0" json:"id"`
	TypeName string `default:"" json:"type_name"`
}
