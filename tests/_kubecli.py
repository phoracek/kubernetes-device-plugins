import urllib3

from kubernetes import config

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

config.load_kube_config('./kubeconfig')
