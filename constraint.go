package main

type Constraint struct {
	ID        uint `json"constraint_id" gorm:"primary_key"`
	Size      int
	GroupID   uint `json:"-"`
	RequestID uint `json:"-"`
}

func (group Group) allConstraints() (*[]Constraint, error) {
	var constraints []Constraint
	result := Db.Model(&group).Association("Constraints").Append(constraints)
	if result.Error != nil {
		return nil, result.Error
	}
	return &constraints, nil
}

func (request Request) allConstraints() (*[]Constraint, error) {
	var constraints []Constraint
	result := Db.Model(&request).Association("Constraints").Append(constraints)
	if result.Error != nil {
		return nil, result.Error
	}
	return &constraints, nil
}
