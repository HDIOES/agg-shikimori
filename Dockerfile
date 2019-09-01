FROM golang
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR $GOPATH/src/github.com/HDIOES/su4na-API-main
COPY Gopkg.toml Gopkg.lock ./
COPY . ./
RUN dep ensure
RUN go install github.com/HDIOES/su4na-API-main
RUN cp configuration-dev.json $GOPATH/bin/
RUN cp -r migrations/ $GOPATH/bin/
WORKDIR $GOPATH/bin
ENTRYPOINT ["./su4na-API-main"]