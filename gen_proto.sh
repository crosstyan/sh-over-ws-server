# https://stackoverflow.com/questions/61666805/correct-format-of-protoc-go-package
protoc --go_out=. --go_opt=paths=source_relative message/*.proto
