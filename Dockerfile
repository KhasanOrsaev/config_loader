FROM golang:1.13
#RUN mkdir /go/src
ARG GIT_LOGIN

ADD . /app
WORKDIR /app

RUN git config --global url."https://git.fin-dev.ru/scm".insteadof "https://git.fin-dev.ru"
RUN go env -w GONOSUMDB="git.fin-dev.ru" && echo "${GIT_LOGIN}" > ~/.netrc && go build -o main
#RUN go build -o main
CMD ["./main"]

EXPOSE 8083
