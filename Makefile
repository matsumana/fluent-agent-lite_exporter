VERSION=$(patsubst "%",%,$(lastword $(shell grep "version\s*=\s" version.go)))
BIN_DIR=bin
BUILD_GOLANG_VERSION=1.8.3
CENTOS_VERSION=7
GITHUB_USERNAME=matsumana

.PHONY : build-with-docker
build-with-docker:
	docker run --rm -v "$(PWD)":/go/src/github.com/matsumana/fluent-agent-lite_exporter -w /go/src/github.com/matsumana/fluent-agent-lite_exporter golang:$(BUILD_GOLANG_VERSION) bash -c 'make build-all'

.PHONY : build-with-circleci
build-with-circleci:
	docker run -v "$(PWD)":/go/src/github.com/matsumana/fluent-agent-lite_exporter -w /go/src/github.com/matsumana/fluent-agent-lite_exporter golang:$(BUILD_GOLANG_VERSION) bash -c 'make build-all'

.PHONY : e2etest-with-circleci
e2etest-with-circleci:
	docker run -v "$(PWD)":/go/src/github.com/matsumana/fluent-agent-lite_exporter -w /go/src/github.com/matsumana/fluent-agent-lite_exporter -e BUILD_GOLANG_VERSION=$(BUILD_GOLANG_VERSION) centos:$(CENTOS_VERSION) bash -c 'yum install -y make && make e2etest'

.PHONY : build-all
build-all: build-linux

.PHONY : build-linux
build-linux:
	make build GOOS=linux GOARCH=amd64

build:
	rm -rf $(BIN_DIR)/fluent-agent-lite_exporter-$(VERSION).$(GOOS)-$(GOARCH)*
	go fmt
	go build -o $(BIN_DIR)/fluent-agent-lite_exporter-$(VERSION).$(GOOS)-$(GOARCH)/fluent-agent-lite_exporter
	tar cvfz $(BIN_DIR)/fluent-agent-lite_exporter-$(VERSION).$(GOOS)-$(GOARCH).tar.gz -C $(BIN_DIR) fluent-agent-lite_exporter-$(VERSION).$(GOOS)-$(GOARCH)

.PHONY : e2etest
e2etest:
	make e2etest_setup
	GOROOT=/usr/local/go GOPATH=/go /usr/local/go/bin/go test -run E2E

.PHONY : e2etest_setup
e2etest_setup:
	# install depend libs
	yum install -y git gcc perl-devel cpan
	curl -L https://cpanmin.us | perl - App::cpanminus

	# install fluent-agent-lite
	git clone https://github.com/tagomoris/fluent-agent-lite.git /root/fluent-agent-lite
	git --git-dir=/root/fluent-agent-lite/.git --work-tree=/root/fluent-agent-lite checkout -b v1.0 refs/tags/v1.0
	/root/fluent-agent-lite/bin/install.sh

	# update config file
	perl -pi -e 's/^# LOGS=\$$\(cat \/etc\/fluent-agent.logs\)/LOGS=\$$\(cat \/etc\/fluent-agent.logs\)/g' /etc/fluent-agent-lite.conf
	echo 'www0  /tmp/www0_access.log' >  /etc/fluent-agent.logs
	echo 'www1  /tmp/www1_access.log' >> /etc/fluent-agent.logs
	echo 'www2  /tmp/www2_access.log' >> /etc/fluent-agent.logs

	# prepare dummy log
	touch /tmp/www0_access.log
	touch /tmp/www1_access.log
	touch /tmp/www2_access.log

	# start
	/etc/init.d/fluent-agent-lite start

	# golang
	yum install -y git
	curl -L https://storage.googleapis.com/golang/go${BUILD_GOLANG_VERSION}.linux-amd64.tar.gz > /tmp/go${BUILD_GOLANG_VERSION}.linux-amd64.tar.gz
	tar xvf /tmp/go${BUILD_GOLANG_VERSION}.linux-amd64.tar.gz -C /usr/local

	# start fluent-agent-lite_exporter
	/go/src/github.com/matsumana/fluent-agent-lite_exporter/bin/fluent-agent-lite_exporter-*.linux-amd64/fluent-agent-lite_exporter -log.level=debug &

	# kill several process for e2etest
	kill `ps aux | grep fluent-agent-lite | grep www0 | awk '{print $$2;}'`
	kill `ps aux | grep fluent-agent-lite | grep www1 | awk '{print $$2;}'`

	# Wait for fluent-agent-lite_exporter to start up
	sleep 3

check-github-token:
	if [ ! -f "./github_token" ]; then echo 'file github_token is required'; exit 1 ; fi

release: build-with-docker check-github-token
	ghr -u $(GITHUB_USERNAME) -t $(shell cat github_token) --draft --replace $(VERSION) $(BIN_DIR)/fluent-agent-lite_exporter-$(VERSION).*.tar.gz
