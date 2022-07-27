// @author cold bin
// @date 2022/7/26

package main

import (
	app "github.com/cold-bin/mini-web"

	"log"
)

type Stu struct {
	Name string `json:"name"`
	Age  int    `json:"-"`
}

func main() {
	engine := app.New()
	engine.POST("/hello", func(c *app.Context) {
		stu := Stu{}
		err := c.ShouldBind(&stu)
		if err != nil {
			log.Println("json: ", err)
			return
		}

		c.JSON(200, stu)

	})

	engine.Run("")
}
