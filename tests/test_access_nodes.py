from . import _cli


def test_connect_to_cluster_node_using_ssh():
    _test_node_hostname('node01')
    _test_node_hostname('node02')


def _test_node_hostname(node):
    rc, stdout, stderr = _cli.run(node, ['hostname'])
    assert rc == 0, stderr
    assert stdout[0] == node, stderr
