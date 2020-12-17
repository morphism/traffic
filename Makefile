all:	test

test:
	go test -v

install:
	cd cmd/traffic && go install
