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
		groups,
	},
	Route{
		"GroupsCheck",
		"POST",
		"/groups/check",
		groupsCheck,
	},
	Route{
		"GroupsNotifications",
		"GET",
		"/groups/{groupId}/notifications",
		groupsNotifications,
	},
	Route{
		"GroupsShow",
		"GET",
		"/groups/{groupId}",
		groupsShow,
	},
	Route{
		"GroupsDelete",
		"DELETE",
		"/groups/{groupId}",
		groupsDelete,
	},
	Route{
		"Requests",
		"GET",
		"/groups/{groupId}/requests",
		requests,
	},
	Route{
		"RequestsCreate",
		"POST",
		"/group/{groupId}/requests{requestId}",
		requestsCreate,
	},
	Route{
		"RequestsDelete",
		"DELETE",
		"/group/{groupId}/requests{requestId}",
		requestsCreate,
	},
	Route{
		"RequestsHistory",
		"GET",
		"/group/{groupId}/requests{requestId}/history",
		requestsHistory,
	},
	Route{
		"Machines",
		"GET",
		"/machines",
		machines,
	},
	Route{
		"MachinesShow",
		"GET",
		"/machines/{machineId}",
		machinesShow,
	},
	Route{
		"MachinesCreate",
		"POST",
		"/machines",
		machinesCreate,
	},
	Route{
		"MachinesDelete",
		"DELETE",
		"/machines/{machineId}",
		machinesDelete,
	},
	Route{
		"Transfer",
		"POST",
		"/tranfer/{file}",
		transfer,
	},
}
