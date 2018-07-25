from kubernetes import client, config

config.load_kube_config('kubeconfig')

v1 = client.CoreV1Api()

print('0')
pod_list = v1.list_namespaced_pod("default")
print('a')
for pod in pod_list.items:
    print("%s\t%s\t%s" % (pod.metadata.name,
                          pod.status.phase,
                          pod.status.pod_ip))
