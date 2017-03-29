package common

import (
	"common/controller/base"
	_ "fmt"
)

type Controller struct {
	base.Controller
}

func (c *Controller) UploadFile() {}

func (c *Controller) DownloadFile() {}

func (c *Controller) AddDocument() {}

func (c *Controller) UpdateDocument() {}
