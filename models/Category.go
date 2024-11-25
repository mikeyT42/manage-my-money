package models

type Category struct {
	Id            int    `db:"id"`
	Description   string `db:"description"`
	Name          string `db:"name"`
	Color         Color  `db:"color"`
	ColorId       int    `db:"color_id"`
	IsUserDefined bool   `db:"is_user_defined"`
}

type Color = string

const (
	None   Color = "none"
	Red    Color = "red"
	Blue   Color = "blue"
	Green  Color = "green"
	Grey   Color = "grey"
	Purple Color = "purple"
	Yellow Color = "yellow"
)
