TEST_FLAGS := -timeout 15m -race -coverpkg=./... ./...

test: export CGO_ENABLED=1
test:
	go get gotest.tools/gotestsum@latest
	gotestsum --junitfile report.xml -- -coverprofile=coverageCI.out $(shell go list ./...) $(TEST_FLAGS)
	go install github.com/t-yuki/gocover-cobertura@latest
	sed 's|$(shell go list -m)/||' coverageCI.out > coverageCI.new; mv coverageCI.new coverageCI.out
	gocover-cobertura < coverageCI.out > coverageCI.xml
