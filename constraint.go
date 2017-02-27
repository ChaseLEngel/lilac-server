package main

import (
	"fmt"
	"strconv"
)

type Constraint struct {
	ID        uint `json:"constraint_id" gorm:"primary_key"`
	Size      int  `json:"size"`
	GroupID   uint `json:"-"`
	RequestID uint `json:"-"`
}

func (group Group) allConstraints() ([]Constraint, error) {
	var constraints []Constraint
	result := Db.Model(&group).Related(&constraints)
	if result.Error != nil {
		return nil, result.Error
	}
	return constraints, nil
}

func (group Group) insertConstraint(constraint *Constraint) error {
	result := Db.Model(&group).Association("Constraints").Append(constraint)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (group Group) findConstraint(id string) (Constraint, error) {
	constraints, err := group.allConstraints()
	if err != nil {
		return Constraint{}, err
	}
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return Constraint{}, err
	}

	for _, constraint := range constraints {
		if constraint.ID == uint(uid) {
			return constraint, nil
		}
	}
	return Constraint{}, fmt.Errorf("record not found")
}

func (group Group) updateConstraint(id string, newCon Constraint) (Constraint, error) {
	constraint, err := group.findConstraint(id)
	if err != nil {
		return Constraint{}, err
	}

	result := Db.Model(&constraint).Updates(newCon)
	if result.Error != nil {
		return Constraint{}, result.Error
	}

	return constraint, nil
}

func (group Group) deleteConstraint(id string) error {
	constraint, err := group.findConstraint(id)
	if err != nil {
		return err
	}

	result := Db.Model(&group).Association("Constraints").Delete(constraint)
	if result.Error != nil {
		return result.Error
	}
	Db.Delete(&constraint)
	return nil
}
