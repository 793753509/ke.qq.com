package web

import "github.com/astaxie/beego"

type Apiserver struct {

}

func NewApiserver() *Apiserver  {
	return &Apiserver{}
}

func (a *Apiserver)Server()  {
	InitRouter()
	beego.Run("127.0.0.1:8099")
}