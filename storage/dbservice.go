package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"ke.qq.com/model"
	"time"
)

type DbOperator struct {
	Operator *sql.DB
}

var inst DbOperator

func setDbOperator(db *sql.DB) {
	inst = DbOperator{Operator: db}
}

func GetDbOperator() DbOperator {
	return inst
}

func InitDbService(mysqlPath string) {
	db, err := sql.Open("mysql", mysqlPath)
	checkErr(err)
	setDbOperator(db)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SyncCourseInfo(courseInfoChan chan model.CourseInfo, tableName string) {
	stmt, err := GetDbOperator().Operator.Prepare(fmt.Sprintf("insert into %s (CourseId, CourseName, Teachers, CourseType, Price) values (?,?,?,?,?)", tableName))
	checkErr(err)
	defer func() {
		err := stmt.Close()
		checkErr(err)
	}()
	for {
		select {
		case info, ok := <-courseInfoChan:
			if !ok {
				return
			}
			err := InsertCurse(info, stmt)
			if err != nil {
				fmt.Printf("insert curse %s error: %s", info.Name, err.Error())
				continue
			}

		}
	}
}

func InsertCurse(info model.CourseInfo, stmt *sql.Stmt) error {
	_, err := stmt.Exec(info.ID, info.Name, info.Teachers, info.CourseType, info.Price)
	if err != nil {
		return err
	}
	return nil
}

type TypeInfo struct {
	Info []TypeCount `json:"信息"`
}

type TypeCount struct {
	Name  string `json:"类别"`
	Count int    `json:"数量"`
}

func QueryTypeCount(tableName string) (TypeInfo, error) {
	var result = TypeInfo{}
	rows, err := GetDbOperator().Operator.Query(fmt.Sprintf("select CourseType,count(CourseType) as TotalCount from %s group by CourseType;", tableName))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var courseType string
		var count int
		err = rows.Scan(&courseType, &count)
		if err != nil {
			return result, err
		}
		result.Info = append(result.Info, TypeCount{
			Name:  courseType,
			Count: count,
		})
	}
	return result, nil
}

func QueryCoursesByType(tableName string, courseType string) ([]*model.CourseInfo, error) {
	result := []*model.CourseInfo{}
	rows, err := GetDbOperator().Operator.Query(fmt.Sprintf("SELECT * FROM %s WHERE CourseType=?;", tableName), courseType)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var courseType string
		var id string
		var name string
		var price string
		var teacher string
		err = rows.Scan(&id, &name, &teacher, &courseType, &price)
		if err != nil {
			return result, err
		}
		c := &model.CourseInfo{
			ID:         id,
			Name:       name,
			Teachers:   teacher,
			CourseType: courseType,
			Price:      price,
		}
		result = append(result, c)
	}
	return result, nil
}

func CreateTable() string {
	tableName := time.Now().Format("2006_01_02")
	sql := "CREATE TABLE IF NOT EXISTS " + tableName +
		"(CourseId VARCHAR(64)," +
		"CourseName VARCHAR(64)," +
		"Teachers LONGTEXT," +
		"CourseType VARCHAR(64)," +
		"Price VARCHAR(64)," +
		"PRIMARY KEY(CourseId),INDEX(CourseType))" +
		"ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	fmt.Println("\n" + sql + "\n")
	stmt, err := GetDbOperator().Operator.Prepare(sql)
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	return tableName

}
