import subprocess


def run(node, command):
    p = subprocess.Popen(
        ['docker', 'exec', 'kubevirt-' + node, 'ssh.sh'] + command,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    stdout, stderr = p.communicate()

    return p.returncode, stdout.split('\n')[:-1], stderr.split('\n')[2:-1]
