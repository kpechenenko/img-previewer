FROM golang:1.23 AS build

ENV BIN_FILE=/opt/image-previewer/image-previewer-app
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} -mod vendor ${CODE_DIR}/cmd/image-previewer/*

FROM alpine:3.9

LABEL SERVICE="image-previewer"

ENV BIN_FILE="/opt/image-previewer/image-previewer-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
