package main

import (
	"fmt"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

func getRoute(name string) (Route, error) {
	for _, route := range routes {
		if route.Name == name {
			return route, nil
		}
	}
	return Route{}, fmt.Errorf("route not found")
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
		"GroupsUpdate",
		"PUT",
		"/groups/{groupId}",
		GroupsUpdate,
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
		"Constraints",
		"GET",
		"/groups/{groupId}/constraints",
		Constraints,
	},
	Route{
		"ConstraintsCreate",
		"POST",
		"/groups/{groupId}/constraints",
		ConstraintsCreate,
	},
	Route{
		"ConstraintsShow",
		"GET",
		"/groups/{groupId}/constraints/{constraintId}",
		ConstraintsShow,
	},
	Route{
		"ConstraintsDelete",
		"DELETE",
		"/groups/{groupId}/constraints/{constraintId}",
		ConstraintsDelete,
	},
	Route{
		"ConstraintsUpdate",
		"PUT",
		"/groups/{groupId}/constraints/{constraintId}",
		ConstraintsUpdate,
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
		"/groups/{groupId}/requests",
		RequestsCreate,
	},
	Route{
		"RequestsDelete",
		"DELETE",
		"/groups/{groupId}/requests/{requestId}",
		RequestsDelete,
	},
	Route{
		"RequestsUpdate",
		"PUT",
		"/groups/{groupId}/requests/{requestId}",
		RequestsUpdate,
	},
	Route{
		"RequestsHistory",
		"GET",
		"/groups/{groupId}/requests/{requestId}/history",
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
