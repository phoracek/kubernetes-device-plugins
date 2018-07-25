from . import _kubectl


def test_list_cluster_nodes():
    rc, stdout, stderr = _kubectl.run(['get', 'nodes'])
    assert rc == 0, stderr
    assert 'node01' in '\n'.join(stdout), stderr
    assert 'node02' in '\n'.join(stdout), stderr
