package controller

import (
	"github.com/vincent-pli/resource-watcher/pkg/controller/resourcewatcher"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, resourcewatcher.Add)
}
