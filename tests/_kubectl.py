import subprocess


def run(command):
    p = subprocess.Popen(
        ['./cluster/kubectl.sh'] + command,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    stdout, stderr = p.communicate()

    return p.returncode, stdout.split('\n')[:-1], stderr.split('\n')[:-1]
