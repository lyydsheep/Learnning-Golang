.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user_gen.go
	@mockgen -source=webook/internal/service/code.go -package=svcmocks -destination=webook/internal/service/mocks/code_gen.go
	@go mod tidy


