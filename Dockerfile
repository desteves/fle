FROM ubuntu:bionic


##########################################################################
# golang https://golang.org/doc/install 
##########################################################################

ENV GOLANG_VERSION_OS_ARCH 1.14.2.linux-amd64
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go && \
    apt-get update && \
    apt-get install -y curl gnupg git make build-essential libssl-dev pkg-config wget && \
    rm -rf /var/lib/apt/lists/*  && \
    curl -sSOL https://dl.google.com/go/go$GOLANG_VERSION_OS_ARCH.tar.gz && \
	tar -C /usr/local -xzf go$GOLANG_VERSION_OS_ARCH.tar.gz
ENV PATH /usr/local/go/bin:/go/bin:$PATH
ENV GOROOT /usr/local/go
ENV GOPATH /go
WORKDIR /go


##########################################################################
# https://github.com/mongodb/libmongocrypt#installing-libmongocrypt-from-distribution-packages
##########################################################################

RUN curl -s https://www.mongodb.org/static/pgp/libmongocrypt.asc | gpg --dearmor >/etc/apt/trusted.gpg.d/libmongocrypt.gpg
RUN echo "deb https://libmongocrypt.s3.amazonaws.com/apt/ubuntu bionic/libmongocrypt/1.0 universe" |  tee /etc/apt/sources.list.d/libmongocrypt.list
RUN apt-get update && apt-get install -y libmongocrypt-dev


##########################################################################
# https://github.com/mongodb-labs/field-level-encryption-sandbox/blob/master/go/golang_fle_install.sh
##########################################################################

RUN  rm -rf /usr/local/lib/libmongocrypt* && \
    /sbin/ldconfig && \
    git clone https://github.com/mongodb/libmongocrypt   && \
    sed -i 's|git@github.com:mongodb|https://github.com/mongodb|g' \
    ./libmongocrypt/.evergreen/prep_c_driver_source.sh  && \
    ./libmongocrypt/.evergreen/compile.sh  && \
    cp -P ./install/libmongocrypt/lib/libmongocrypt* /usr/local/lib && \
    /sbin/ldconfig  && \
    /sbin/ldconfig -p | grep crypt && \
    rm -rf libmongocrypt  mongo-c-driver install 


##########################################################################
# https://docs.mongodb.com/manual/reference/security-client-side-encryption-appendix/#installation
##########################################################################

RUN wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -
RUN echo "deb [ arch=amd64 ] http://repo.mongodb.com/apt/ubuntu bionic/mongodb-enterprise/4.2 multiverse" |  tee /etc/apt/sources.list.d/mongodb-enterprise.list
RUN apt-get update && apt-get install mongodb-enterprise-cryptd


##########################################################################
# FLE Demo
##########################################################################

EXPOSE  8888 27020
WORKDIR /go/src/github.com/desteves/
RUN git clone https://github.com/desteves/fle
WORKDIR /go/src/github.com/desteves/fle
RUN go build -tags cse main.go
ENTRYPOINT [ "./main" ]