FROM golang:1.15.2 AS builder

WORKDIR /dist
COPY . .
RUN go build -i .


FROM ubuntu:20.04 AS runtime
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y lsb-release wget software-properties-common g++
RUN bash -c "$(wget -O - https://apt.llvm.org/llvm.sh)"
RUN ln -s /usr/bin/clang++-11 /usr/bin/clang++

EXPOSE 8080
WORKDIR /dist
COPY --from=builder /dist/inout_tester /dist/inout_tester
