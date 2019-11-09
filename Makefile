make-zip:
	GOOS=linux go build smartagri.go
	zip function.zip smartagri
