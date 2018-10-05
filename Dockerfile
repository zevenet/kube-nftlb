# Download the latest Debian image
FROM debian:latest

# Default shell when executing RUN
SHELL ["/bin/bash", "-c"]

# Install required tools
RUN apt-get update
RUN apt-get install -y \
    git \
    bison \
    flex \
    binutils \
    build-essential \
    autoconf \
    libtool \
    pkg-config \
    libgmp-dev \
    libreadline-dev \
    libjansson-dev \
    libev-dev \
    cmake \
    curl \
    dnsutils

# Download the most recent version of libmnl and install it
RUN git clone git://git.netfilter.org/libmnl/
WORKDIR /libmnl
RUN sh autogen.sh
RUN ./configure
RUN make
RUN make install
WORKDIR /

# Download the most recent version of libnftnl and install it 
RUN git clone git://git.netfilter.org/libnftnl
WORKDIR /libnftnl
RUN sh autogen.sh
RUN ./configure
RUN make
RUN make install
WORKDIR /

# Download the most recent versión of nftables and install it
RUN ldconfig
RUN git clone git://git.netfilter.org/nftables
WORKDIR /nftables
RUN sh autogen.sh
RUN ./configure
RUN make
RUN make install
WORKDIR /
RUN ldconfig

# Download the most recent version of nftlb and install it
RUN git clone https://github.com/zevenet/nftlb
WORKDIR /nftlb
RUN autoreconf -fi
RUN ./configure
RUN make
RUN make install
WORKDIR /
RUN ldconfig

# For testing: after compiling the Golang client, copy the binary to / and make it the entrypoint
COPY ./app /app
ENTRYPOINT /app
