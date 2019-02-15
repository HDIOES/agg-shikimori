FROM golang
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get -v github.com/rubenv/sql-migrate/...
WORKDIR $GOPATH/src/github.com/HDIOES/cpa-backend
COPY Gopkg.toml Gopkg.lock ./
COPY . ./
RUN dep ensure
RUN go install github.com/HDIOES/cpa-backend
#here we need to apply sql migrations
RUN cp configuration.json $GOPATH/bin/
WORKDIR $GOPATH/bin

