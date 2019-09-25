FROM ubuntu:bionic AS offline

RUN apt-get update && apt-get install -y git curl

ENV BUCC_VERSION 9977aad
RUN git clone https://github.com/starkandwayne/bucc /bucc-build && \
cd /bucc-build && git checkout ${BUCC_VERSION}

RUN echo '{"internal_ip":"10.1.1.10"}' > /bucc-build/vars.yml
RUN /bucc-build/bin/bucc offline --cpi docker --destination /tmp/offline
RUN mkdir -p /bucc && tar -xf /tmp/offline/bucc-*.tar.gz -C /bucc

FROM ubuntu:bionic
RUN apt-get update && apt-get install -y build-essential zlibc zlib1g-dev \
    ruby ruby-dev openssl libxslt1-dev libxml2-dev libssl-dev libreadline7 \
    libreadline-dev libyaml-dev libsqlite3-dev sqlite3 jq curl \
    net-tools iputils-ping

COPY --from=offline /bucc/bucc-build /bucc
RUN cd bucc && ln -sf state/vars.yml && \
    echo "source <(/bucc/bin/bucc env)" >> ~/.bashrc

VOLUME /bucc/state

ENTRYPOINT /bucc/bin/bucc
