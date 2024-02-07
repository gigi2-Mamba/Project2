.PHONY: mock
mock:
	@mockgen -source=./internal/service/article.go -package=svcmocks -destination=./internal/service/mocks/article-mock0.go
	@mockgen -source=./internal/service/code.go -package=svcmocks  -destination=.internal/service/mocks/user-code.go
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=.internal/service/mocks/user-user.go
	@mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=.internal/service/sms/mocks/sms-mock.go
	@mockgen -source=./internal/repository/user.go -package=reposmocks -destination=.internal/repository/mocks/user-mock.go
	@mockgen -source=./internal/repository/code.go -package=reposmocks -destination=.internal/repository/mocks/code-mock.go
	@mockgen -source=./internal/repository/article.go -package=reposmocks -destination=.internal/repository/mocks/article-mock.go
	@mockgen -source=./internal/repository/article_author.go -package=reposmocks -destination=.internal/repository/mocks/article_author.go
	@mockgen -source=./internal/repository/article_reader.go -package=reposmocks -destination=.internal/repository/mocks/article_reader.go
	@mockgen -source=./internal/repository/cache/user.go -package=cachemocks -destination=.internal/repository/cache/mocks/cache-user-mock.go
	@mockgen -source=./internal/repository/dao/user.go -package=daomocks -destination=.internal/repository/dao/mocks/cache-dao-mock.go
	@mockgen -source=./internal/repository/dao/article_author.go -package=daomocks -destination=.internal/repository/dao/mocks/article_author.go
	@mockgen -source=./internal/repository/dao/article_reader.go -package=daomocks -destination=.internal/repository/dao/mocks/article_reader.go
	@mockgen -source=./pkg/limiter/types.go -package=limitermocks -destination=./pkg/limiter/mocks/litmiter-mock.go
	@go mod tidy

