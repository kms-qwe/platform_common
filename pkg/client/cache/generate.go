package cache

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i RedisCache  -o ./mocks/ -s "_minimock.go"
