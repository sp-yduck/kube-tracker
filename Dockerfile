FROM golang:1.19-alpine3.17 as gobuilder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/ ./...


FROM bash:alpine3.16
RUN apk update && apk add curl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.26.0/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl
RUN apk add openssh git
# RUN echo -e "Host github.com\n\tStrictHostKeyChecking no\n" >> ~/.ssh/config
COPY --from=gobuilder /usr/local/bin/ /usr/local/bin/
WORKDIR /root/
COPY entrypoint.sh /root/entrypoint.sh
# ENTRYPOINT [ "/root/enrypoint.sh" ]