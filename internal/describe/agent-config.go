/*
 *  *******************************************************************************
 *  * Copyright (c) 2020 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package describe

import (
	"fmt"
	"strings"

	"github.com/eclipse-iofog/iofog-go-sdk/v2/pkg/client"
	rsc "github.com/eclipse-iofog/iofogctl/v2/internal/resource"
	iutil "github.com/eclipse-iofog/iofogctl/v2/internal/util"

	"github.com/eclipse-iofog/iofogctl/v2/internal/config"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/iofog"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/util"
)

type agentConfigExecutor struct {
	namespace string
	name      string
	filename  string
}

func newAgentConfigExecutor(namespace, name, filename string) *agentConfigExecutor {
	a := &agentConfigExecutor{}
	a.namespace = namespace
	a.name = name
	a.filename = filename
	return a
}

func (exe *agentConfigExecutor) GetName() string {
	return exe.name
}

func getAgentNameFromUUID(agentMapByUUID map[string]client.AgentInfo, uuid string) (name string) {
	if uuid == iofog.VanillaRouterAgentName {
		return uuid
	}
	agent, found := agentMapByUUID[uuid]
	if !found {
		util.PrintNotify(fmt.Sprintf("Could not find router uuid %s\n", uuid))
		name = "UNKNOWN ROUTER: " + uuid
	} else {
		name = agent.Name
	}
	return
}

func (exe *agentConfigExecutor) Execute() error {
	ns, err := config.GetNamespace(exe.namespace)
	if err != nil {
		return err
	}
	// Get config
	agent, err := ns.GetAgent(exe.name)
	if err != nil {
		return err
	}

	// Connect to controller
	ctrl, err := iutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	// Get all agents for mapping uuid to name if required
	getAgentList, err := ctrl.ListAgents(client.ListAgentsRequest{})
	if err != nil {
		return err
	}
	// Map by uuid for easier access
	agentMapByUUID := make(map[string]client.AgentInfo)
	for _, agent := range getAgentList.Agents {
		agentMapByUUID[agent.UUID] = agent
	}

	getAgentResponse, err := ctrl.GetAgentByID(agent.GetUUID())
	if err != nil {
		// The agents might not be provisioned with Controller
		if strings.Contains(err.Error(), "NotFoundError") {
			return util.NewInputError("Cannot describe an Agent that is not provisioned with the Controller in namespace " + exe.namespace)
		}
		return err
	}

	fogType, found := rsc.FogTypeIntMap[getAgentResponse.FogType]
	if !found {
		fogType = "auto"
	}

	routerConfig := client.RouterConfig{
		RouterMode:      &getAgentResponse.RouterMode,
		MessagingPort:   getAgentResponse.MessagingPort,
		EdgeRouterPort:  getAgentResponse.EdgeRouterPort,
		InterRouterPort: getAgentResponse.InterRouterPort,
	}

	var upstreamRoutersPtr *[]string

	if getAgentResponse.UpstreamRouters != nil {
		upstreamRouters := []string{}
		for _, upstreamRouterAgentUUID := range *getAgentResponse.UpstreamRouters {
			upstreamRouters = append(upstreamRouters, getAgentNameFromUUID(agentMapByUUID, upstreamRouterAgentUUID))
		}
		upstreamRoutersPtr = &upstreamRouters
	}

	var networkRouterPtr *string
	if getAgentResponse.NetworkRouter != nil {
		networkRouter := getAgentNameFromUUID(agentMapByUUID, *getAgentResponse.NetworkRouter)
		networkRouterPtr = &networkRouter
	}

	agentConfig := rsc.AgentConfiguration{
		Name:        getAgentResponse.Name,
		Location:    getAgentResponse.Location,
		Latitude:    getAgentResponse.Latitude,
		Longitude:   getAgentResponse.Longitude,
		Description: getAgentResponse.Description,
		FogType:     &fogType,
		AgentConfiguration: client.AgentConfiguration{
			DockerURL:                 &getAgentResponse.DockerURL,
			DiskLimit:                 &getAgentResponse.DiskLimit,
			DiskDirectory:             &getAgentResponse.DiskDirectory,
			MemoryLimit:               &getAgentResponse.MemoryLimit,
			CPULimit:                  &getAgentResponse.CPULimit,
			LogLimit:                  &getAgentResponse.LogLimit,
			LogDirectory:              &getAgentResponse.LogDirectory,
			LogFileCount:              &getAgentResponse.LogFileCount,
			StatusFrequency:           &getAgentResponse.StatusFrequency,
			ChangeFrequency:           &getAgentResponse.ChangeFrequency,
			DeviceScanFrequency:       &getAgentResponse.DeviceScanFrequency,
			BluetoothEnabled:          &getAgentResponse.BluetoothEnabled,
			WatchdogEnabled:           &getAgentResponse.WatchdogEnabled,
			AbstractedHardwareEnabled: &getAgentResponse.AbstractedHardwareEnabled,
			LogLevel:                  getAgentResponse.LogLevel,
			DockerPruningFrequency:    getAgentResponse.DockerPruningFrequency,
			AvailableDiskThreshold:    getAgentResponse.AvailableDiskThreshold,
			UpstreamRouters:           upstreamRoutersPtr,
			NetworkRouter:             networkRouterPtr,
			RouterConfig:              routerConfig,
		},
	}

	header := config.Header{
		APIVersion: config.LatestAPIVersion,
		Kind:       config.AgentConfigKind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      exe.name,
		},
		Spec: agentConfig,
	}

	if exe.filename == "" {
		if err = util.Print(header); err != nil {
			return err
		}
	} else {
		if err = util.FPrint(header, exe.filename); err != nil {
			return err
		}
	}
	return nil
}
