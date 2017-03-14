package main

import (
	"strconv"
	"testing"
)

func compareGroups(g1, g2 Group) bool {
	if g1.Name != g2.Name {
		return false
	} else if g1.Link != g2.Link {
		return false
	} else if g1.DownloadPath != g2.DownloadPath {
		return false
	}
	return true
}

func newMockGroup() *Group {
	group := new(Group)
	group.Name = "Test Group"
	group.Link = "http://localhost/feed"
	group.DownloadPath = "/home/user/files"
	return group
}

func formatUint(u uint) string {
	u64 := uint64(u)
	return strconv.FormatUint(u64, 10)
}

func TestInsertGroup(t *testing.T) {
	setup()
	defer teardown()

	group := newMockGroup()
	err := group.insert()
	if err != nil {
		t.Fatal(err)
	}

	var groups []Group
	result := Db.Find(&groups)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	if len(groups) != 1 {
		t.Error("Not only record in table.")
	}

	if !compareGroups(*group, groups[0]) {
		t.Errorf("\nExpected %v\nActual %v", *group, groups[0])
	}
}

func TestFindGroup(t *testing.T) {
	setup()
	defer teardown()
	group := newMockGroup()
	result := Db.Create(&group)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	foundGroup, err := findGroup(formatUint(group.ID))
	if err != nil {
		t.Fatal(err)
	}

	if !compareGroups(*group, foundGroup) {
		t.Errorf("\nExpected %v\nActual %v", *group, foundGroup)
	}
}

func TestDeleteGroup(t *testing.T) {
	setup()
	defer teardown()
	group := newMockGroup()
	result := Db.Create(&group)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	deleteGroup(formatUint(group.ID))
	var groups []Group
	if err := Db.Find(&groups).Error; err != nil {
		t.Fatal(err)
	}

	if len(groups) != 0 {
		t.Error("Group was not deleted:\n%v", groups)
	}
}

func TestUpdateLastChecked(t *testing.T) {
	setup()
	defer teardown()
	group := newMockGroup()
	initalTime := group.LastChecked
	err := group.updateLastChecked()
	if err != nil {
		t.Fatal(err)
	}
	if initalTime == group.LastChecked {
		t.Errorf("Last Checked was not updated:\nintial time:%v\nafter:%v\n", initalTime, group.LastChecked)
	}
}

func TestAllGroups(t *testing.T) {
	setup()
	defer teardown()

	var expected []Group

	group1 := newMockGroup()
	result := Db.Create(&group1)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	expected = append(expected, *group1)
	group2 := newMockGroup()
	result = Db.Create(&group2)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	expected = append(expected, *group2)

	groups, err := allGroups()
	if err != nil {
		t.Fatal(err)
	}

	for index, group := range groups {
		if !compareGroups(expected[index], group) {
			t.Errorf("\nExpected %v\nActual %v\n", expected[index], group)
		}
	}
}

func TestUpdateGroup(t *testing.T) {
	setup()
	defer teardown()
	group := newMockGroup()
	result := Db.Create(&group)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	updatedGroup := newMockGroup()
	updatedGroup.Name = "Updated name"
	updatedGroup.Link = "Updated link"
	updatedGroup.DownloadPath = "Updated download path"

	group.update(*updatedGroup)

	if !compareGroups(*group, *updatedGroup) {
		t.Errorf("\nExpected %v\nActual %v\n", updatedGroup, group)
	}
}
