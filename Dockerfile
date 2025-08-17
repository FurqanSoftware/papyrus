FROM ubuntu:noble-20250529

MAINTAINER Mahmud Ridwan <ridwan@furqansoftware.com>

RUN apt update && \
    DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
