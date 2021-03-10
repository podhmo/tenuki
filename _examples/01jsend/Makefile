default:
	$(MAKE) -k test

test: test-text test-json test-openapi

test-text:
	go clean . && CAPTURE_WRITEFILE=tetdata/text go test
test-json:
	go clean . && CAPTURE_WRITEFILE=tetdata/json CAPTURE_LAYOUT=json go test
test-openapi:
	go clean . && CAPTURE_WRITEFILE=tetdata/openapi CAPTURE_LAYOUT=openapi go test
