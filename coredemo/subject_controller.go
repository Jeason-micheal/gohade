package main

import (
	"github.com/gohade/hade/framework/gin"
)

func SubjectAddController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectAddController")
}

func SubjectListController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectListController")
}

func SubjectDelController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectDeleteController")
}

func SubjectUpdateController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectUpdateController")
}

func SubjectGetController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectGetController")
}

func SubjectNameController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectNameController")
}
