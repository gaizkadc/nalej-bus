FROM alpine:edge

# gmock package is only available on the edge alpine version

RUN apk add --no-cache --virtual build-dependencies \
    ca-certificates \
    git \
    cmake \
    alpine-sdk \
    zlib-dev \
    boost-dev \
    boost-static \
    protobuf-dev \
    curl-dev \
    gtest-dev \
    && git clone https://github.com/apache/pulsar.git \
    && cd pulsar \
    && git checkout tags/v2.3.1 \
    && cd pulsar-client-cpp \
    && cmake . -DBUILD_TESTS=OFF -DBUILD_PYTHON_WRAPPER=OFF -DLINK_STATIC=OFF \
    && make -j 2 \
    && make install \
    && cd ../.. \
    && rm -rf pulsar \
    && rm -rf .github \
    && apk -v del build-dependencies \
    && rm -rf /var/cache/apk/

    # the library is now available at /usr/local/lib/libpulsar.so
    # dev headers are available at /usr/local/include/pulsar

    # If BUILD_PYTHON_WRAPPER is enabled
    # python-dev

    # If BUILD_TESTS is enabled
    # gtest
    # gmock

    # IF LINK_STATIC is enabled
    # boost-static

    # The static linking process requires all libraries to be compiled with the -fPIC flag. It seems that
    # the boost-static libraries are not prepared to use this flag.