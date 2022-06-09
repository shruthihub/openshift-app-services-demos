# **ACK RDS w/ SBO Annotations Demo**

## Info

Example demo built off of [AWS ACK RDS](https://github.com/aws-controllers-k8s/rds-controller) using SBO Annotations to bind to an application. This demo also utilizes [minikube](https://minikube.sigs.k8s.io/docs/start/),  [Docker](https://docs.docker.com/get-docker/) and [helm](https://helm.sh/).

## About

This demo accesses API functionality of [AWS RDS](https://aws.amazon.com/rds/?trk=c0fcea17-fb6a-4c27-ad98-192318a276ff&sc_channel=ps&sc_campaign=acquisition&sc_medium=ACQ-P|PS-GO|Brand|Desktop|SU|Database|Solution|US|EN|Text&s_kwcid=AL!4422!3!548665196304!e!!g!!amazon%20relational%20db&ef_id=EAIaIQobChMIpJKl5a-C9wIVbfHjBx0QYAOwEAAYASABEgJe7PD_BwE:G:s&s_kwcid=AL!4422!3!548665196304!e!!g!!amazon%20relational%20db) using the RDS controller from [Amazon Controllers for Kubernetes](https://github.com/aws-controllers-k8s/community).  Furthermore, this demo utilizes [Service Binding Operator](https://github.com/redhat-developer/service-binding-operator) (SBO) and CR annotations to automatically bind a DBInstance to an example application. 

## Setup
### AWS ACK RDS
[AWS ACK RDS](https://github.com/redhat-developer/sbo-services-investigation/tree/main/aws/s3/app) must be installed on the cluster. Refer to the AWS ACK [installation guide](https://aws-controllers-k8s.github.io/community/docs/user-docs/install/).

### rds-postgre-chart-demo
This is the helm chart that sets up and deploys a DBInstance on the cluster. In the `values.yaml` file, change any values of the DBInstance as needed. The `dbcreds.username` & `dbcreds.password` are undefined and must be set prior to deploying the helm chart. These values can be whatever the user wants, and they will serve as the `Master Username` and `Master Password`. 

Inside `templates/dbinstance-postgres`, the SBO annotations can be found underneath `metadata.annotations`. These annotations are preconfigured for the example application given in the demo, but feel free to change these annotations to bind data specific to the application you are using. More info on how to correctly define SBO annotations can be found [here](https://redhat-developer.github.io/service-binding-operator/userguide/exposing-binding-data/adding-annotation.html).

### rds-app-chart-demo
This helm chart deploys an CRUD application that binds to the DBInstance using the SBO object and annotations found on the `dbinstance-postgres.yaml` file. However, any postgres application can be used as substitute. If using a different postgres application, the annotations in `dbinstance-postgres.yaml` should be changes according to the binding info the application needs. 

### app-files
This folder contains the source code and Dockerfile to build the docker image for the example application. If using a custom postgres application, ignore this step. Otherwise, create the docker image in the directory `../app-files` with the commands:

`eval $(minikube docker-env)`
`docker build -t aws-rds-test-api .`


### Custom Application
If a custom postgres application is being used in place of the example application provided, the user must edit their `sb-rds-endpoint-demo.yaml` SBO object. Underneath `.spec.application`, the `resource`, `group`, and `version` must correspond to the application deployment. 

Additionally, in the custom app deployment, a label under `.metadata.labels` must be created and match the label found in `sb-rds-endpoint-demo.yaml`.

## Installation
### DBInstance
1. Using helm in the `../demo` directory of the project, run the command:

	`helm install rds-postgre-chart-demo -n [namespace] rds-postgre-chart-demo`

	where `[namespace]` is the namespace where ACK RDS is installed on.

2. Run the command to very the DBInstance has been installed:

	`kubectl get dbinstances -n [namespace]`

3. Run this command using the name of the DBInstance created:

	`kubectl describe dbinstance [dbinstance-name] -n [namespace]`

4. The field `.status.DBInstanceStatus` will initially be in the state of `creating`. Before continuing the DBInstance must reach the state of `available`. The initialization process may take some time.

5. Once the DBInstance is in the state of `.status.DBInstanceStatus.available`, the user can then apply their application to the cluster. If using a custom application, apply it to the cluster. If using the included example application, continue following the instructions bellow.

### Application
1. If the docker image has been made for the example application, then deploy the app chart using the command: 

	`helm install rds-app-chart-demo -n [namespace] rds-app-chart-demo`

2. Due to exponential back-off of SBO object when looking for applications to bind to, the application may not bind immediately. This status of the binding can be seen using the command:

	`kubectl get servicebindings -n [namespace]`

3.  Once the application is bound, the binding can be confirmed grabbing the working pod and using the command:

	`kubectl get pods -n [namespace]`
	`kubectl exec [pod] -n [namespace] -- env`
	
	The bound variables can be found with the prefix `DBINSTANCE_` inside the pod's environment. At this point, a custom application is successfully bound and good to use with the AWS RDS DBInstance created. If using the example application, follow the steps bellow to use the application.

### Example Application Walkthrough

1. Once the DBInstance is available, a table needs to be inserted using an external program (such as pgAdmin). This table is to be named `myTable`, with the columns `uid:integer` & `test:text`. This setup is needed for the functionality of our CRUD application.

2. Once the table is setup, the application can be serviced. Run the commands:
	
	`kubectl expose deployment aws-rds-sbo-demo -n [namespace] --type=NodePort --name=rds-demo-svc`
	
	`minikube service -n [namespace] rds-demo-svc`

3. The database table can now be viewed through the address provided by minikube. 

	`GET` calls on the landing page displays the contents of `myTable`

	`POST` calls on the address `../insert` will post a table entry to the application
	
	`DELETE` calls on the address `../delete/?uid=#` will delete the table entry with the specified `uid`


	
