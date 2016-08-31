from subprocess import Popen, PIPE
import sys

def handler(event=None, context=None):
    p = Popen(['./ecr-cleaner', "-aws.region", "eu-west-1", "-dry-run", "true"], stdin=PIPE, stdout=PIPE, stderr=PIPE, bufsize=-1)

    while True:
        out = p.stderr.read(1)
        if out == '' and p.poll() != None:
            break
        if out != '':
            sys.stdout.write(out)
            sys.stdout.flush()

    return p.returncode