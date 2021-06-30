default: cli
		@echo -e '\t[CB-Spider] build ./bin/cb-spider....'
		@go mod download
		@go build -o bin/cb-spider ./api-runtime
dyna plugin plug dynamic: cli
		@echo -e '\t[CB-Spider] build ./bin/cb-spider with plugin mode...'
		@go mod download
	        @go build -tags dyna -o bin/cb-spider-dyna ./api-runtime
		@./build_all_driver_lib.sh;
cc:
		@echo -e '\t[CB-Spider] build ./bin/cb-spider-arm for arm...'
	        GOOS=linux GOARCH=arm go build -o cb-spider-arm ./api-runtime
clean clear:
		@echo -e '\t[CB-Spider] cleaning...'
	        @rm -rf bin/cb-spider bin/cb-spider-dyna bin/cb-spider-arm
	        @rm -rf dist-tmp

cli-dist: cli
		@echo -e '\t[CB-Spider] tar spctl... to dist'
		@mkdir -p /tmp/spider/dist/conf 
		@cp ./interface/spctl ./interface/grpc_conf.yaml /tmp/spider/dist 1> /dev/null
		@cp ./conf/log_conf.yaml /tmp/spider/dist/conf 1> /dev/null
		@mkdir -p ./dist
		@tar -zcvf ./dist/spctl-`(date +%Y.%m.%d.%H)`.tar.gz -C /tmp/spider/dist ./ 1> /dev/null
		@rm -rf /tmp/spider
cli:
		@echo -e '\t[CB-Spider] build ./interface/spctl...'
		@go mod download
		@go build -o ./interface/spctl ./interface/cli/spider/spider.go
swag swagger:
		@echo -e '\t[CB-Spider] build Swagger docs'
		@~/go/bin/swag i -g api-runtime/rest-runtime/CBSpiderRuntime.go -o api-runtime/rest-runtime/docs

