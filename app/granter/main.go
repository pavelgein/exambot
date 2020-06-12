package main

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mkideal/cli"
	"github.com/pavelgein/exambot/internal/db"
	"github.com/pavelgein/exambot/internal/models"
	"github.com/pavelgein/exambot/internal/oauth"
)

type UserOptions struct {
	cli.Helper
	Name        string `cli:"*n,name" usage:"user name"`
	Token       string `cli:"*t,token" usage:"user token"`
	IsSuperUser bool   `cli:"superuser" usage:"make user superuser"`
}

type GrantOptions struct {
	cli.Helper
	Name string `cli:"*n,name" usage:"user name"`
	Page string `cli:"*p,page" usage:"page"`
}

func CreateUser(checker *oauth.OAuthMultiPageChecker, options *UserOptions) error {
	checker.CreateUser(options.Name, options.Token, options.IsSuperUser)
	return nil
}

func Grant(checker *oauth.OAuthMultiPageChecker, options *GrantOptions) error {
	user, err := checker.GetUser(options.Name)
	if err != nil {
		return err
	}

	page, err := checker.GetPage(options.Page)
	if err != nil {
		return err
	}

	return checker.GrantPermission(&user, &page)
}

func main() {
	db, err := db.CreateDBFromEnvironment()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.ApiUser{}, &models.Page{}, &models.Role{})
	db.AutoMigrate(&models.Task{}, &models.Assignment{}, &models.Course{}, &models.User{}, &models.TelegramUser{}, &models.TaskSet{})

	checker := oauth.OAuthMultiPageChecker{
		DB:   db,
		Salt: os.Getenv("EXAMBOT_SALT"),
	}

	createUserCommand := &cli.Command{
		Name: "user",
		Desc: "create user with given token",
		Argv: func() interface{} { return new(UserOptions) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*UserOptions)
			return CreateUser(&checker, argv)
		},
	}

	grantCommand := &cli.Command{
		Name: "grant",
		Desc: "grant permission to user",
		Argv: func() interface{} { return new(GrantOptions) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*GrantOptions)
			return Grant(&checker, argv)
		},
	}

	helpCommand := cli.HelpCommand("show help")
	rootCommand := &cli.Command{
		Name: "root",
		Desc: "Manage users and permissions",
		Fn: func(*cli.Context) error {
			return nil
		},
	}

	if err := cli.Root(
		rootCommand,
		cli.Tree(createUserCommand),
		cli.Tree(helpCommand),
		cli.Tree(grantCommand),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
