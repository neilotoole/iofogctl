/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package cmd

import (
	"fmt"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
	"github.com/spf13/cobra"
)

// Set by linker
var (
	version   = "undefined"
	commit    = "undefined"
	buildDate = "undefined"
	platform  = "undefined"
)

type versionSpec struct {
	Version   string `yaml:"version"`
	Commit    string `yaml:"commit"`
	BuildDate string `yaml:"buildDate"`
	Platform  string `yaml:"platform"`
}

func newVersionCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get CLI application version",
		Run: func(cmd *cobra.Command, args []string) {
			spec := versionSpec{
				Version:   version,
				Commit:    commit,
				BuildDate: buildDate,
				Platform:  platform,
			}
			fmt.Printf("\033[38;5;117mCopyright (C) 2019, Edgeworx, Inc.\033[0m\n")
			util.Print(spec)
		},
	}
	return cmd
}
