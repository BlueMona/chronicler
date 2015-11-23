package timeline

import (
	riak "github.com/basho/riak-go-client"
	gosnow "github.com/sdming/gosnow"
)

var snowFlake *gosnow.SnowFlake
