package qgomysql

import (
	"testing"
)

func TestQueryParser(test *testing.T) {

	db := Config{}
	db.Prepare("select * from users where id = :id and name = :name")
	db.BindParam("id", 100)
	db.BindParam("name", "Test")
	result := db.ToString()

	if result != "select * from users where id = 100 and name = \"Test\"" {
		test.Error("TestQueryParser fail")
	}
}

func TestServerParser(test *testing.T) {

	db := Config{}
	db.SERVER.Delimiter = "$$$$"

	urls := map[string]string{
		"/select/users/id/100/name/Test/$$$$/1/":                               "select * from users where id = \"100\" and name = \"Test\" limit 1",
		"/delete/users/id/100/name/Test/$$$$/1/":                               "delete from users where id = \"100\" and name = \"Test\" limit 1",
		"/insert/users/id/100/name/Test/":                                      "insert into users set id = \"100\", name = \"Test\"",
		"/update/users/id/9/parent_id/100/$$$$/name/Test/lastnane/SecondTest/": "update users set name = \"Test\", lastnane = \"SecondTest\" where id = \"9\" and parent_id = \"100\"",
	}

	var result string

	for path, expect := range urls {

		db.buildQuery(path)
		result = db.ToString()

		if result != expect {
			test.Error("TestServerParser fail")
		}
	}

}
