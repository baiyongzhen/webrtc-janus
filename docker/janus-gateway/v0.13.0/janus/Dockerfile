FROM debian:latest

ENV BUILD_FOLDER /tmp/build
ENV WORK_FOLDER /opt

RUN mkdir -p ${BUILD_FOLDER}
RUN mkdir -p ${WORK_FOLDER}

RUN apt-get update && apt-get install -y --fix-missing \
  libmicrohttpd-dev \
  libjansson-dev \
  libnice-dev \
  libssl-dev \
  libsofia-sip-ua-dev \
  libglib2.0-dev \
  libopus-dev \
  libogg-dev \
  libcurl4-openssl-dev \
  liblua5.3-dev \
  libconfig-dev \
  pkg-config \
  gengetopt \
  libtool \
  automake \
  gtk-doc-tools \
  doxygen \
  graphviz \
  libnanomsg-dev \
  libunwind-dev \
  cmake \
  golang-go

RUN apt-get install -y wget \
  && apt-get install -y vim \
  && apt-get install -y git \
  && apt-get install -y curl

RUN PATH="/usr/local/bin:$PATH"
RUN yes | apt-get install python3-pip
RUN yes | pip3 install meson
RUN yes | pip3 install ninja

RUN cd ${BUILD_FOLDER} \
  && git clone https://gitlab.freedesktop.org/libnice/libnice --depth=1 \
  && cd ./libnice \
  && meson --prefix=/usr --buildtype=release build \
  && ninja -C build \
  && ninja -C build install

RUN cd ${BUILD_FOLDER} \
  && wget https://github.com/cisco/libsrtp/archive/v2.3.0.tar.gz \
  && tar xfv v2.3.0.tar.gz \
  && cd libsrtp-2.3.0 \
  && ./configure --prefix=/usr --enable-openssl \
  && make shared_library \
  && make install

RUN cd ${BUILD_FOLDER} \
  && git clone https://github.com/sctplab/usrsctp -b 0.9.5.0 --depth=1 \
  && cd ./usrsctp \
  && ./bootstrap \
  && ./configure --prefix=/usr --libdir=/usr/lib \
  && make \
  && make install

RUN cd ${BUILD_FOLDER} \
  && git clone https://libwebsockets.org/repo/libwebsockets -b v3.2-stable --depth=1 \
  && cd ./libwebsockets \
  && mkdir build \
  && cd ./build \
  && cmake -DCMAKE_BUILD_TYPE=Release \
        -DLWS_MAX_SMP=1 -DLWS_WITHOUT_EXTENSIONS=0 \
        -DCMAKE_INSTALL_PREFIX:PATH=/usr \
        -DCMAKE_C_FLAGS="-fpic" .. \
  && make \
  && make install

RUN cd ${WORK_FOLDER} \
  && git clone https://github.com/meetecho/janus-gateway.git -b v0.13.0 --depth=1 \
  && cd janus-gateway \
  && sh autogen.sh \
  && ./configure --prefix=/opt/janus --enable-docs \
  && make \
  && make install \
  && make configs


  CMD ["/opt/janus/bin/janus"]