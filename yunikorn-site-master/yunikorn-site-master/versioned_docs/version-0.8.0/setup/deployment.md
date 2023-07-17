---
id: deployment
title: Deployment Guide
---

<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->

The easiest way to deploy YuniKorn is to leverage our [helm charts](https://github.com/apache/incubator-yunikorn-k8shim/tree/master/helm-charts),
you can find the guide [here](get_started/user_guide.md). However, if you want to explore the deployment procedure
step by step, here are the instructions.

## Setup RBAC

The first step is to create the RBAC role for the scheduler, see [yunikorn-rbac.yaml](https://github.com/apache/incubator-yunikorn-k8shim/blob/master/deployments/scheduler/yunikorn-rbac.yaml)
```
kubectl create -f scheduler/yunikorn-rbac.yaml
```
The role is a requirement on the current versions of kubernetes.

## Create the ConfigMap

YuniKorn loads its configuration from a K8s configmap, so it is required to create the configmap before launching the scheduler.

- download a sample configuration file:
```
curl -o queues.yaml https://raw.githubusercontent.com/apache/incubator-yunikorn-k8shim/master/conf/queues.yaml
```
- create ConfigMap in kubernetes:
```
kubectl create configmap yunikorn-configs --from-file=queues.yaml
```
- check if the ConfigMap was created correctly:
```
kubectl describe configmaps yunikorn-configs
```

For more information about how to manage scheduler's configuration via configmap, see more from [here](setup/configure_scheduler.md).

## Deploy the scheduler on k8s

Before you can deploy the scheduler the image for the scheduler and the web interface must be build with the appropriate tags.
The procedure on how to build the images is explained in the [build document](get_started/developer_guide.md). See [scheduler.yaml](https://github.com/apache/incubator-yunikorn-k8shim/blob/master/deployments/scheduler/scheduler.yaml)
```
kubectl create -f scheduler/scheduler.yaml
```
The deployment will run 2 containers from your pre-built docker images in 1 pod,

* yunikorn-scheduler-core (yunikorn scheduler core and shim for K8s)
* yunikorn-scheduler-web (web UI)

The pod is deployed as a customized scheduler, it will take the responsibility to schedule pods which explicitly specifies `schedulerName: yunikorn` in pod's spec.

## Access to the web UI

When the scheduler is deployed, the web UI is also deployed in a container.
Port forwarding for the web interface on the standard ports can be turned on via:

```
POD=`kubectl get pod -l app=yunikorn -o jsonpath="{.items[0].metadata.name}"` && \
kubectl port-forward ${POD} 9889 9080
```

`9889` is the default port for Web UI, `9080` is the default port of scheduler's Restful service where web UI retrieves info from.
Once this is done, web UI will be available at: http://localhost:9889.
