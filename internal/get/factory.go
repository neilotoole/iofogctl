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

package get

import (
	"github.com/eclipse-iofog/iofogctl/v2/internal/execute"
	"github.com/eclipse-iofog/iofogctl/v2/pkg/util"
)

func NewExecutor(resourceType, namespace string, showDetached bool) (execute.Executor, error) {
	switch resourceType {
	case "namespaces":
		return newNamespaceExecutor(), nil
	case "all":
		return newAllExecutor(namespace, showDetached), nil
	case "controllers":
		return newControllerExecutor(namespace), nil
	case "agents":
		return newAgentExecutor(namespace, showDetached), nil
	case "microservices":
		return newMicroserviceExecutor(namespace), nil
	case "applications":
		return newApplicationExecutor(namespace), nil
	case "catalog":
		return newCatalogExecutor(namespace), nil
	case "registries":
		return newRegistryExecutor(namespace), nil
	case "volumes":
		return newVolumeExecutor(namespace), nil
	case "routes":
		return newRouteExecutor(namespace), nil
	default:
		msg := "Unknown resource: '" + resourceType + "'"
		return nil, util.NewInputError(msg)
	}
}
