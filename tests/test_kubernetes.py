from kubernetes import client

from . import _kubectl
from . import _kubecli  # NOQA: initialize kubernetes client

import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def test_list_cluster_nodes_using_script():
    rc, stdout, stderr = _kubectl.run(['get', 'nodes'])
    assert rc == 0, stderr
    assert 'node01' in '\n'.join(stdout), stderr
    assert 'node02' in '\n'.join(stdout), stderr


def test_list_cluster_nodes_using_python_client():
    v1 = client.CoreV1Api()
    node_list = v1.list_node()
    node_name_list = [node.metadata.name for node in node_list.items]
    assert set(['node01', 'node02']) == set(node_name_list)
