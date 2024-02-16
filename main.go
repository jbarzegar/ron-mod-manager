/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/jbarzegar/ron-mod-manager/cmd"
	stateManagement "github.com/jbarzegar/ron-mod-manager/state-management"
)

func main() {
	stateManagement.PreflightChecks()
	cmd.Execute()
}
