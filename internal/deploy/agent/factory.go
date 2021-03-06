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

package deployagent

import (
	"fmt"

	"github.com/eclipse-iofog/iofogctl/v2/internal/config"
	"github.com/eclipse-iofog/iofogctl/v2/internal/execute"
	rsc "github.com/eclipse-iofog/iofogctl/v2/internal/resource"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/iofog/install"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/util"
)

type AgentDeployExecutor interface {
	execute.Executor
	GetHost() string
}

type facadeExecutor struct {
	isSystem  bool
	exe       execute.Executor
	agent     rsc.Agent
	namespace string
}

func (facade *facadeExecutor) GetHost() string {
	return facade.agent.GetHost()
}

func (facade *facadeExecutor) Execute() (err error) {
	// Check the namespace exists
	ns, err := config.GetNamespace(facade.namespace)
	if err != nil {
		return err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return err
	}

	// Check Controller exists
	if len(controlPlane.GetControllers()) == 0 {
		return util.NewInputError("This namespace does not have a Controller. You must first deploy a Controller before deploying Agents")
	}

	if !facade.isSystem || install.IsVerbose() {
		util.SpinStart(fmt.Sprintf("Deploying agent %s", facade.GetName()))
	}

	if err = facade.exe.Execute(); err != nil {
		return
	}

	// Don't add system agent to the namespace config file
	if !facade.isSystem {
		if err = ns.UpdateAgent(facade.agent); err != nil {
			return
		}
	}
	return config.Flush()
}

func (facade *facadeExecutor) GetName() string {
	return facade.exe.GetName()
}

func (facade *facadeExecutor) ProvisionAgent() (string, error) {
	// Required for attach
	provisionExecutor, ok := facade.exe.(execute.ProvisioningExecutor)
	if !ok {
		return "", util.NewInternalError("Facade executor: Could not convert executor")
	}
	return provisionExecutor.ProvisionAgent()
}

func newFacadeExecutor(exe execute.Executor, namespace string, agent rsc.Agent, isSystem bool) execute.Executor {
	return &facadeExecutor{
		exe:       exe,
		namespace: namespace,
		isSystem:  isSystem,
		agent:     agent,
	}
}

func NewRemoteExecutor(namespace string, agent *rsc.RemoteAgent, isSystem bool) (execute.Executor, error) {
	if err := util.IsLowerAlphanumeric("Agent", agent.GetName()); err != nil {
		return nil, err
	}

	if err := agent.Sanitize(); err != nil {
		return nil, err
	}
	if err := agent.ValidateSSH(); err != nil {
		return nil, err
	}
	return newFacadeExecutor(newRemoteExecutor(namespace, agent), namespace, agent, isSystem), nil
}

func NewLocalExecutor(namespace string, agent *rsc.LocalAgent, isSystem bool) (execute.Executor, error) {
	if err := util.IsLowerAlphanumeric("Agent", agent.GetName()); err != nil {
		return nil, err
	}

	exe, err := newLocalExecutor(namespace, agent, isSystem)
	if err != nil {
		return nil, err
	}
	return newFacadeExecutor(exe, namespace, agent, isSystem), nil
}
