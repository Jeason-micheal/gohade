package middleware

import (
	"fmt"
	"github.com/gohade/hade/framework/gin"
)

func Test1() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
	}
}

func Test2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
		panic("Test2 panic test.....")
	}

}

func Test3() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test3")
		c.Next()
		fmt.Println("middleware post test3")
	}
}
