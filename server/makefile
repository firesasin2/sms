
VERSION=0.1
RELEASE=0.1.beta
BUILDTIME=(2021-04-16)

LDFLAGS="-X main.Version=${VERSION} -X main.Build=${RELEASE}${BUILDTIME} -w -s"
export LDFLAGS

all :
	go build -ldflags ${LDFLAGS} -o server

clean :
	@rm -f ./server