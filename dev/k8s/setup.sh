#!/bin/sh
go run ../../cmd/server/main.go create member -t one.org
go run ../../cmd/server/main.go create member -t two.org
go run ../../cmd/server/main.go create relationship -a one.org -b two.org
go run ../../cmd/server/main.go generate token -t one.org                
go run ../../cmd/server/main.go generate token -t two.org