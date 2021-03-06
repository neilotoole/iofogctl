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

package pruneagent

import (
	"fmt"
	"strings"

	rsc "github.com/eclipse-iofog/iofogctl/v2/internal/resource"
	iutil "github.com/eclipse-iofog/iofogctl/v2/internal/util"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/iofog/install"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/util"
)

func (exe executor) remoteAgentPrune(agent rsc.Agent) error {
	ctrl, err := iutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}
	// If controller exists, prune the agent
	// Perform Docker pruning of Agent through Controller
	if err = ctrl.PruneAgent(agent.GetUUID()); err != nil {
		if !strings.Contains(err.Error(), "NotFoundError") {
			return err
		}
	}
	return nil
}

func (exe executor) remoteDetachedAgentPrune(agent *rsc.RemoteAgent) error {
	if err := agent.ValidateSSH(); err != nil {
		return err
	} else {
		sshAgent := install.NewRemoteAgent(agent.SSH.User, agent.Host, agent.SSH.Port, agent.SSH.KeyFile, agent.Name, agent.UUID)
		if err := sshAgent.Prune(); err != nil {
			return util.NewInternalError(fmt.Sprintf("Failed to Prune Iofog resource %s. %s", agent.Name, err.Error()))
		}
	}
	return nil
}
