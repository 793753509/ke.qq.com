package web

import (
	"github.com/astaxie/beego"
	"ke.qq.com/storage"
	"net/http"
	"net/url"
)

type ErrResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CourseTypeController struct {
	beego.Controller
}

func (q *CourseTypeController) Get() {
	q.Controller.EnableRender = false
	tableName := q.Ctx.Input.Param(":data")
	info, err := storage.QueryTypeCount(tableName)
	if err != nil {
		q.Ctx.Output.Status = http.StatusBadRequest
		q.Data["json"] = &ErrResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		q.ServeJSON()
		return
	}
	q.Data["json"] = info
	q.ServeJSON()
}

type CourseInfoController struct {
	beego.Controller
}

func (q *CourseInfoController) Get() {
	q.Controller.EnableRender = false
	tableName := q.Ctx.Input.Param(":data")
	courseType, err := url.QueryUnescape(q.Ctx.Input.Param(":type"))

	if err != nil {
		q.Ctx.Output.Status = http.StatusBadRequest
		q.Data["json"] = &ErrResult{
			Code:    http.StatusBadRequest,
			Message: "course type unescape error",
		}
		q.ServeJSON()
		return
	}

	info, err := storage.QueryCoursesByType(tableName, courseType)
	if err != nil {
		q.Ctx.Output.Status = http.StatusBadRequest
		q.Data["json"] = &ErrResult{
			Code:    http.StatusBadRequest,
			Message: "query courses error: " + err.Error(),
		}
		q.ServeJSON()
		return
	}
	q.Data["json"] = info
	q.ServeJSON()
}
