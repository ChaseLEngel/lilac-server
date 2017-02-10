package main

import (
	"github.com/chaselengel/lilac/group"
	"github.com/chaselengel/lilac/machine"
	"github.com/chaselengel/lilac/request"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Login",
		"POST",
		"/login",
		login,
	},
	Route{
		"Groups",
		"GET",
		"/groups",
		groups.Groups,
	},
	Route{
		"GroupsCheck",
		"POST",
		"/groups/check",
		groups.GroupsCheck,
	},
	Route{
		"GroupsNotifications",
		"GET",
		"/groups/{groupId}/notifications",
		groups.GroupsNotifications,
	},
	Route{
		"GroupsShow",
		"GET",
		"/groups/{groupId}",
		groups.GroupsShow,
	},
	Route{
		"GroupsDelete",
		"DELETE",
		"/groups/{groupId}",
		groups.GroupsDelete,
	},
	Route{
		"Requests",
		"GET",
		"/groups/{groupId}/requests",
		request.Requests,
	},
	Route{
		"RequestsCreate",
		"POST",
		"/group/{groupId}/requests{requestId}",
		request.RequestsCreate,
	},
	Route{
		"RequestsDelete",
		"DELETE",
		"/group/{groupId}/requests{requestId}",
		request.RequestsCreate,
	},
	Route{
		"RequestsHistory",
		"GET",
		"/group/{groupId}/requests{requestId}/history",
		request.RequestsHistory,
	},
	Route{
		"Machines",
		"GET",
		"/machines",
		machine.Machines,
	},
	Route{
		"MachinesShow",
		"GET",
		"/machines/{machineId}",
		machine.MachinesShow,
	},
	Route{
		"MachinesCreate",
		"POST",
		"/machines",
		machine.MachinesCreate,
	},
	Route{
		"MachinesDelete",
		"DELETE",
		"/machines/{machineId}",
		machine.MachinesDelete,
	},
	Route{
		"Transfer",
		"POST",
		"/tranfer/{file}",
		transfer,
	},
}
