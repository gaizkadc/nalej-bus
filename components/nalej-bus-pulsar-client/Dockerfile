FROM alpine:edge

# gmock package is only available on the edge alpine version

RUN apk add --no-cache --virtual build-dependencies \
    ca-certificates \
    git \
    cmake \
    alpine-sdk \
    zlib-dev \
    boost-dev \
    python-dev \
    protobuf-dev \
    curl-dev \
    gtest-dev \
    gtest \
    gmock \
    && rm -rf /var/cache/apk/ \
    && git clone https://github.com/apache/pulsar.git \
    && cd pulsar \
    && git checkout tags/v2.3.1 \
    && cd pulsar-client-cpp \
    && cmake . \
    && make -j 2 \
    && make install \
    && cd ../.. \
    && rm -rf pulsar \
    && apk -v del build-dependencies
    # library is now available at /usr/local/lib/libpulsar.so
    # headers are available at /usr/local/include/pulsar