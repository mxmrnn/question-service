FROM ubuntu:latest
LABEL authors="max_m"

ENTRYPOINT ["top", "-b"]