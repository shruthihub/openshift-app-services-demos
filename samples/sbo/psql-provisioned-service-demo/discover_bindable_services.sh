#!/usr/bin/env bash

sbo_object=$(oc get servicebinding.servicebinding.io -o json)

sbo_spec_list=$(echo $sbo_object | jq -rc '.items[].spec | @base64')
echo Bindable Services:
for sbo_spec in $sbo_spec_list; do
    echo "-------------------------------------"
    label_selector=$(echo $sbo_spec | base64 -D | jq -rc .workload.selector.matchLabels)
    service_name=$(echo $sbo_spec | base64 -D | jq -r .service.name)
    service_type=$(echo $sbo_spec | base64 -D | jq -r .service.kind)
    if [[ $label_selector != "null" ]]; then
        echo "Label Selector: $label_selector"
        echo "Service Name: $service_name"
        echo "Service Type: $service_type"
        secret_name=$(oc get $service_type $service_name -o json | jq -r .status.binding.name)
        echo "Service Endpoint Definition: $secret_name"
        echo "Service Endpoint Definition Fields:"
        oc get secret $secret_name -o json | jq -r '.data | to_entries[] | .key'
    fi
done