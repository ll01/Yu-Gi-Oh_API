package main

type TableNames struct {
}

func (tableNames *TableNames) Archtype() string {
	return "archtype_table"
}
func (tableNames *TableNames) Attribute() string {
	return "attribute_table"
}
func (tableNames *TableNames) EffectKeyword() string {
	return "effect_keyword_table"
}
func (tableNames *TableNames) ForeignName() string {
	return "foreign_name_table"
}
func (tableNames *TableNames) LinkArrow() string {
	return "LinkArrow_table"
}
