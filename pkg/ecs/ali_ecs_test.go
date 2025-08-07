package ecs

import (
	"fmt"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

func TestCreateAliEcsClient(t *testing.T) {
	ecsClient, err := CreateEcsClient()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if ecsClient != nil {
		fmt.Println("Ecs connection successful")
	} else {
		fmt.Println("Ecs connection failure")
	}
}
func TestCreateAliEcsTry(t *testing.T) {
	ecsClient, err := CreateEcsClient()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if ecsClient != nil {
		fmt.Println("Ecs connection successful")
	} else {
		fmt.Println("Ecs connection failure")
		return
	}

	name := "os-dazz-gateway"

	res, err := CreateEcsForTenant(ecsClient, "test.com", name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*res)

}

func TestGetEcsInfo(t *testing.T) {
	ecsClient, err := CreateEcsClient()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if ecsClient != nil {
		fmt.Println("Ecs connection successful")
	} else {
		fmt.Println("Ecs connection failure")
		return
	}
	res, err := DescribeInstanceStatus(ecsClient, []*string{
		tea.String("i-uf6fjqdjzjt036sbgguo"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*res)
}
