FROM golang:1.6-onbuild

RUN make clean && make
ENTRYPOINT ["/go/src/app/bin/linux/ecr-cleaner"]
