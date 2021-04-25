package web

import "github.com/astaxie/beego"

func InitRouter(){
	beego.Router("/api/v1/coursetype/:data",&CourseTypeController{})
	beego.Router("/api/v1/courseinfo/:data/:type",&CourseInfoController{})
}
