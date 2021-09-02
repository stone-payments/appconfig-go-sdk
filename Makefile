TIME_SHORT	= `date +%H:%M:%S`
TIME		= $(TIME_SHORT)
CNone        := $(shell printf "\033[0m")
GREEN        := $(shell printf "\033[32m")
RED          := $(shell printf "\033[31m")

OK		= echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL	= (echo ${TIME} ${RED}[FAIL]${CNone} && false)

test:
	go test -v ./... || $(FAIL)
	@$(OK) tests passed

coverage:
	go test -v -coverprofile=cover.out ./...
	go tool cover -html=cover.out
