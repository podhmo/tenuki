SED := $(shell which gsed || which sed)

sync:
	rm -f dump.go
	wget https://raw.githubusercontent.com/golang/go/master/src/net/http/httputil/dump.go
	$(SED) -i 's/NewChunkedWriter(/httputil.NewChunkedWriter(/g' dump.go
	$(SED) -i 's/io\.Discard/ioutil.Discard/g' dump.go
	$(SED) -i 's/io\.NopCloser/ioutil.NopCloser/g' dump.go
	goimports -w dump.go
.PHONY: sync
