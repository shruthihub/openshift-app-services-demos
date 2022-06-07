This sample uses the [ssm](https://github.com/dperaza4dustbit/ssm) application which requires a Posgresql Database to store some data. The Posgresql database will come from the [PSQL Operator](https://github.com/dperaza4dustbit/psql-operator), which support the [Service Binding Specification](https://github.com/servicebinding/spec). The main goal of this sample is to show how SBO can be used keeping in mind the admin and the developer roles. 

### Requirement:
SBO Operator installed on OpenShift 4.10 via OLM
PSQLOperator installed:
1. Clone https://github.com/dperaza4dustbit/psql-operator
1. Run make install to configure the new CRD and corresponding structures
1. Run make deploy to deploy the PSQL controller

### Admin Actions:
helm install psql-service-claim ../chart/psql-service-claim

### Developer Actions:
./discover_bindable_services.sh
helm show values ../chart/ssm-service 
cat ssm_values.yaml
helm install davpssm ../chart/ssm-service -f ssm_values.yaml
oc get servicebinding.servicebinding.io psql-service-ssm-app
kubectl port-forward --namespace default svc/ssm-usa 8080:8080
./test_ssm.sh