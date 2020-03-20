package controller

import (
	"znc-operator/pkg/controller/znc"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, znc.Add)
}
