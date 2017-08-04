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
		"/groups/{groupId}/check",
		GroupsCheck,
	},
	Route{
		"GroupsUpdate",
		"PUT",
		"/groups/{groupId}",
		GroupsUpdate,
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
		"Settings",
		"GET",
		"/groups/{groupId}/settings",
		Settings,
	},
	Route{
		"SettingsCreate",
		"POST",
		"/groups/{groupId}/settings",
		SettingsCreate,
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
		"RequestsHistoryDelete",
		"DELETE",
		"/groups/{groupId}/requests/{requestId}/history/{historyId}",
		RequestsHistoryDelete,
	},
	Route{
		"Machines",
		"GET",
		"/machines",
		Machines,
	},
	Route{
		"MachinesCreate",
		"POST",
		"/machines",
		MachinesCreate,
	},
	Route{
		"MachinesUpdate",
		"PUT",
		"/machines/{machineId}",
		MachinesUpdate,
	},
	Route{
		"MachinesDelete",
		"DELETE",
		"/machines/{machineId}",
		MachinesDelete,
	},
	Route{
		"RequestMachine",
		"GET",
		"/requests/{requestID}/machines",
		RequestMachines,
	},
	Route{
		"RequestMachineCreate",
		"POST",
		"/requests/{requestID}/machines",
		RequestMachinesCreate,
	},
	Route{
		"TransferRequest",
		"POST",
		"/transfer",
		transferRequest,
	},
}
