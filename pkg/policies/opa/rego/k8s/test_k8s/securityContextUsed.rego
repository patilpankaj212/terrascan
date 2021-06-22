package accurics

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
    some i;
    pod := input.kubernetes_pod[_]
    container := pod.config.spec.containers[i]
    not container.securityContext
    traverse := sprintf("spec.containers[%d].securityContext",[i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}

{{.prefix}}{{.name}}{{.suffix}}[retVal] {
    some i;
    pod := input.kubernetes_pod[_]
    initcontainer := pod.config.spec.initContainers[i]
    not initcontainer.securityContext
    traverse := sprintf("spec.initContainers[%d].securityContext",[i])
    retVal := {"Id": pod.id, "Traverse": traverse}
}