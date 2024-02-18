FROM golang:alpine AS builder
ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/exec ./main.go

FROM ubuntu:latest as runner

#Install dependencies
RUN apt update
RUN apt install -y sqlite3 unzip wget
RUN apt install -y lib32gcc-s1
RUN apt install -y libsdl2-2.0-0
RUN apt install -y libgdiplus

#Install steamcmd
WORKDIR /steamcmd
RUN wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz && \
    tar -xvzf steamcmd_linux.tar.gz && \
    chmod +x /steamcmd/steamcmd.sh && \
    chmod +x /steamcmd/linux32/steamcmd && \
    chmod +x /steamcmd/linux32/steamerrorreporter && \
    rm steamcmd_linux.tar.gz

#run build with following args, to force build container from this line
#docker-compose build rust-server --build-arg server_version=$(date +%s)
ARG server_version=unknown

#Download server
RUN /bin/bash ./steamcmd.sh +force_install_dir /server +login anonymous +app_update 258550 validate +quit +quit

#Download oxide
WORKDIR /server
RUN wget -P /server https://github.com/OxideMod/Oxide.Rust/releases/latest/download/Oxide.Rust-linux.zip
RUN unzip -o Oxide.Rust-linux.zip
RUN rm Oxide.Rust-linux.zip


COPY --from=builder /app/exec /server/wrapper
RUN chmod +x /server/RustDedicated
RUN chmod +x /server/wrapper

CMD ["/server/wrapper"]