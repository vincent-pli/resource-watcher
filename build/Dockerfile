FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/resource-watcher \
    USER_UID=1001 \
    USER_NAME=resource-watcher

# install operator binary
COPY build/_output/bin/resource-watcher ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
