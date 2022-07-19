APP_REPO ?= quay.io/pmacik/go-rest-mongodb
APP_IMAGE ?= $(APP_REPO):latest

NAMESPACE ?= default

MULTIARCH_PLATFORMS ?= linux/amd64,linux/arm64,linux/s390x,linux/ppc64le

KUBERNETES_RUNTIME ?= openshift

RUNTIME_VERSION ?= default

CLI ?= kubectl -n $(NAMESPACE)

.PHONY: namespace
namespace:
ifeq ($(KUBERNETES_RUNTIME), openshift)
	-oc new-project $(NAMESPACE)
else ifeq ($(KUBERNETES_RUNTIME), minikube)
	-kubectl create namespace $(NAMESPACE)
	kubectl config set-context --current --namespace=$(NAMESPACE)
endif

.PHONY: build-app
build-app:
	docker build -f Dockerfile -t $(APP_IMAGE) .

.PHONY: push-app
push-app:
	docker push $(APP_IMAGE)

.PHONY: deploy-app
deploy-app:
	$(CLI) apply -f app-deployment.$(KUBERNETES_RUNTIME).$(RUNTIME_VERSION).yaml

.PHONY: deploy-mongodb
deploy-mongodb:
	$(CLI) apply -f mongo-cluster.yaml
	$(CLI) wait --for=condition=ready=True PerconaServerMongoDB/mongo-cluster --timeout=5m

.PHONY: mongo-uri
mongo-uri:
	@echo mongodb://$$(kubectl get secret mongo-cluster-secrets -o jsonpath='{.data.MONGODB_USER_ADMIN_USER}' | base64 -d):$$(kubectl get secret mongo-cluster-secrets -o jsonpath='{.data.MONGODB_USER_ADMIN_PASSWORD}' | base64 -d)@$$(kubectl get PerconaServerMongoDB mongo-cluster -o jsonpath='{.status.host}'):27017/admin?ssl=false


.PHONY: bind-mongodb
bind-mongodb:
	$(CLI) apply -f go-rest-mongodb-binding.yaml

.PHONY: unbind-mongodb
unbind-mongodb:
	$(CLI) delete -f go-rest-mongodb-binding.yaml

