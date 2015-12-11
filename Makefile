prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-temporal; then rm -rf src/github.com/whosonfirst/go-whosonfirst-temporal; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-temporal
	cp temporal.go src/github.com/whosonfirst/go-whosonfirst-temporal/