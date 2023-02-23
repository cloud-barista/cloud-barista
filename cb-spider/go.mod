module github.com/cloud-barista/cb-spider

go 1.19

replace (
	github.com/IBM/vpc-go-sdk/0.23.0 => github.com/IBM/vpc-go-sdk v0.23.0
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/docker/distribution => github.com/docker/distribution v2.8.0+incompatible
)

retract (
	v0.3.12
	v0.3.11
)
