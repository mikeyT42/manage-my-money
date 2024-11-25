package models

type CategoryColor struct {
	Id    int   `db:"id"`
	Color Color `db:"color"`
}
