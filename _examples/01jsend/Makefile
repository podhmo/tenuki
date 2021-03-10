default:
	rm -r testdata/{text,json,openapi}
	$(MAKE) -k test

test: test-text test-json test-openapi

test-text:
	go clean . && CAPTURE_WRITEFILE=testdata/text go test
test-json:
	go clean . && CAPTURE_WRITEFILE=testdata/json CAPTURE_LAYOUT=json go test
test-openapi:
	go clean . && CAPTURE_WRITEFILE=testdata/openapi CAPTURE_LAYOUT=openapi go test
