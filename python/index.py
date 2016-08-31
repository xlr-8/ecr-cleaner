from subprocess import Popen, PIPE
import sys

def handler(event=None, context=None):
    cmd = ['./ecr-cleaner'] + parseJsonToParams(event)
    p = Popen(cmd, stdin=PIPE, stdout=PIPE, stderr=PIPE, bufsize=-1)

    while True:
        out = p.stderr.read(1)
        if out == '' and p.poll() != None:
            break
        if out != '':
            sys.stdout.write(out)
            sys.stdout.flush()

    return p.returncode


def parseJsonToParams(jsonEvent):
    params = []

    for param in jsonEvent:
        params.append(param["name"])
        params.append(param["value"])

    return params