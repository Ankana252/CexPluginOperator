# CexPluginOperatorOpenshift

* Operator watches config map on kubernetes cex plugin
* The cex plugin is designed to run on s390 architecture with RHEL
* Changes were made on the plugin code base to make it run on minikube with Ubuntu
* The main function of the plugin is simple an infinite loop with some sleep duration
* The edited cex plugin code can be found in the fork: https://github.com/Ankana252/k8s-cex-dev-plugin/tree/minikube_experimentation
Steps to compile Operator SDK on WSL (Makefile of SDK is meant to run on Linux)
* make install
* make docker-build IMG=cex-operator:local
* make deploy IMG=cex-operator:local

Run operator locally (easiest for testing)
* Ensure Docker is pointing to the minikube cluster
* Ensure CRDs are installed: kubectl apply -f config/crd/bases/
* Run operator: go run cmd/main.go

Test
* In another terminal, edit the config map: kubectl edit configmap cex-resources -n cex-device-plugin
* Watch operator logs in the first terminal to see reconciliation trigger