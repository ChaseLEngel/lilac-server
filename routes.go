package main

import (
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
		Groups,
	},
	Route{
		"GroupsCheck",
		"POST",
		"/groups/check",
		GroupsCheck,
	},
	Route{
		"GroupsNotifications",
		"GET",
		"/groups/{groupId}/notifications",
		GroupsNotifications,
	},
	Route{
		"GroupsCreate",
		"POST",
		"/groups",
		GroupsCreate,
	},
	Route{
		"GroupsShow",
		"GET",
		"/groups/{groupId}",
		GroupsShow,
	},
	Route{
		"GroupsDelete",
		"DELETE",
		"/groups/{groupId}",
		GroupsDelete,
	},
	Route{
		"Requests",
		"GET",
		"/groups/{groupId}/requests",
		Requests,
	},
	Route{
		"RequestsCreate",
		"POST",
		"/group/{groupId}/requests{requestId}",
		RequestsCreate,
	},
	Route{
		"RequestsDelete",
		"DELETE",
		"/group/{groupId}/requests{requestId}",
		RequestsCreate,
	},
	Route{
		"RequestsHistory",
		"GET",
		"/group/{groupId}/requests{requestId}/history",
		RequestsHistory,
	},
	Route{
		"Machines",
		"GET",
		"/machines",
		Machines,
	},
	Route{
		"MachinesShow",
		"GET",
		"/machines/{machineId}",
		MachinesShow,
	},
	Route{
		"MachinesCreate",
		"POST",
		"/machines",
		MachinesCreate,
	},
	Route{
		"MachinesDelete",
		"DELETE",
		"/machines/{machineId}",
		MachinesDelete,
	},
	Route{
		"Transfer",
		"POST",
		"/tranfer/{file}",
		transfer,
	},
}
