test:
	go test
.PHONY: test

test-with-writefile:
	CAPTURE_WRITEFILE=writefile go test
.POHNY: test-with-writefile

clean:
	rm -rf writefile

