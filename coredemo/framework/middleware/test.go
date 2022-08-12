package middleware

import (
	"fmt"
	"github.com/gohade/hade/framework"
)

func Test1() framework.ControllerHandler {

	return func(c *framework.Context) error {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
		return nil
	}
}

func Test2() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
		panic("Test2 panic test.....")
		return nil
	}

}

func Test3() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
		return nil
	}
}
