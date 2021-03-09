test:
	go test -v ./...
	go test ./_examples/*
.PHONY: test

test-with-writefile:
	CAPTURE_WRITEFILE=writefile go test -v -count=1
.POHNY: test-with-writefile

clean:
	rm -rf writefile

