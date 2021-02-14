package integration_tests

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	fmt.Println("create db")
	fmt.Println("create rabbit")
	fmt.Println("create app")
	fmt.Println("create scheduler")
	fmt.Println("create sender")

	fmt.Println("init app")

	fmt.Println("check create event")
	fmt.Println("check get events")

	fmt.Println("add events sleep 6 sec")
	fmt.Println("get event from rabbitmq")

	fmt.Println("shutdows all containers")
}
