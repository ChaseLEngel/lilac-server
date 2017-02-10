package groups

import (
	"github.com/chaselengel/lilac/constraint"
	"github.com/chaselengel/lilac/notification"
	"github.com/chaselengel/lilac/request"
)

type Group struct {
	DownloadPath  string
	Link          string
	Request       []request.Request
	History       []request.MatchHistory
	Constraints   []constraint.Constraint
	Notifications []notification.Notification
}
