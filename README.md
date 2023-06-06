# Mysql request

## Config
```go
	db := &db.Config{
		SSH: db.SSHconfig{
			Use_ssh:  true,
			User:     "ssh_user",
			Password: "ssh_passwword",
			Port:     22,
			Address:  "0.0.0.0",
		},
		MYSQL: db.MYSQLconfig{
			User:     "mysql_user",
			Password: "mysql_password",
		},
		SERVER: db.SERVERconfig{
			Port:      5555,
			Delimiter: "@@@",
		},
	}
```

## Use in code
```go
	db := &db.Config{
		MYSQL: db.MYSQLconfig{
			User:     "mysql_user",
			Password: "mysql_password",
		},
	}

	db.Prepare("select * from db.users where parent_id = :parent_id and name = :name ")
	db.BindParam("parent_id", 3)
	db.BindParam("name", "User")
	result := db.Execute()

	fmt.Println(result)
```

## Server
```go
	db := &db.Config{
		MYSQL: db.MYSQLconfig{
			User:     "mysql_user",
			Password: "mysql_password",
		},
		SERVER: db.SERVERconfig{
			Port:      8080,
			Delimiter: "$$$$",
		},
	}

	db.RunServer()
```

<sub>insert into db.users set parent_id = 3, name = "User"</sub>
> address:port/insert/db.users/parent_id/3/name/User/

<sub>select * from db.users where parent_id = 3 and name = "User" limit 1</sub>
> address:port/select/db.users/parent_id/3/name/User/$$$$/1/

<sub>update db.users set parent_id = 5, name = "NewUser" where parent_id = 3 and name = "User"</sub>
> address:port/update/db.users/parent_id/3/name/User/$$$$/parent_id/5/name/NewUser/

<sub>delete from db.users where parent_id = 5 and name = "NewUser" limit 1</sub>
> address:port/delete/db.users/parent_id/5/name/NewUser/$$$$/1/

