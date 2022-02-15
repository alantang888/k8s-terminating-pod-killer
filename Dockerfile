FROM golang:1.17 as build

WORKDIR /go/src/github.com/alantang888/k8s-terminating-pod-killer
COPY . .
ENV GO111MODULE=on
RUN go mod download
WORKDIR /go/src/github.com/alantang888/k8s-terminating-pod-killer/cmd/k8s-terminating-pod-killer
RUN go build -o /go/bin/app


FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /

CMD ["/app"]
