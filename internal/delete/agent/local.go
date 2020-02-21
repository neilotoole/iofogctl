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

package deleteagent

import (
	"fmt"
	"strings"

	"github.com/eclipse-iofog/iofogctl/pkg/util"

	"github.com/eclipse-iofog/iofogctl/pkg/iofog/install"
)

func (exe executor) deleteLocalContainer() error {
	client, err := install.NewLocalContainerClient()
	if err != nil {
		return err
	}

	// Clean agent containers (normal and system)
	if errClean := client.CleanContainer(install.GetLocalContainerName("agent", false)); errClean != nil {
		util.PrintNotify(fmt.Sprintf("Could not clean Agent container: %v", errClean))
	}
	// if errClean := client.CleanContainer(install.GetLocalContainerName("agent", true)); errClean != nil {
	// 	util.PrintNotify(fmt.Sprintf("Could not clean Agent container: %v", errClean))
	// }

	// Clean microservices
	containers, err := client.ListContainers()
	for _, container := range containers {
		fmt.Printf("Container names: %v\n", container.Names)
		for _, containerName := range container.Names {
			if strings.HasPrefix(containerName, "/iofog_") {
				fmt.Printf("Deleting name: %s\n", containerName)
				if errClean := client.CleanContainerByID(container.ID); errClean != nil {
					util.PrintNotify(fmt.Sprintf("Could not clean Microservice container: %v", errClean))
				}
			}
		}
	}

	return nil
}
