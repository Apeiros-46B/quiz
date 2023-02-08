package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

// {{{ get default public dir
func defaultPublicDir() string {
    if strings.HasPrefix(os.Args[0], os.TempDir()) {
        return "./pb_public"
    }
    return filepath.Join(os.Args[0], "../pb_public")
}
// }}}

// {{{ check submitted answers to quiz questions
func checkSubmissions(dao *daos.Dao, submissions map[string]int) bool {
    correctAnswers, err := dao.FindRecordsByExpr("correct")
    if err != nil { panic(err) }

    score := 0
    for _, v := range correctAnswers {
        question := v.Get("question")
        correct := v.GetInt("correct")
        if submissions[question.(string)] == correct {
            score++
        }
    }

    return score == len(correctAnswers)
}
// }}}

var uIDRegex = regexp.MustCompile("[A-Za-z0-9_\\-]{21}")
var numRegex = regexp.MustCompile("[0-9]+")

func check(pattern *regexp.Regexp, s string) bool {
    return pattern.FindString(s) != ""
}

func main() {
    app := pocketbase.New()

    var publicDirFlag string

    // add "--publicDir" flag
    app.RootCmd.PersistentFlags().StringVar(
        &publicDirFlag,
        "publicDir",
        defaultPublicDir(),
        "the directory to serve static files",
    )

    // {{{ setup migrations
    migrationsDir := "" // default to "pb_migrations" (for js) and "migrations" (for go)

    // load js files to allow loading external JavaScript migrations
    jsvm.MustRegisterMigrations(app, &jsvm.MigrationsOptions{
        Dir: migrationsDir,
    })

    // register the `migrate` command
    migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
        TemplateLang: migratecmd.TemplateLangJS, // or migratecmd.TemplateLangGo (default)
        Dir:          migrationsDir,
        Automigrate:  false,
    })
    // }}}

    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        // serve static files from pub dir
        e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDirFlag), true))

        // {{{ api endpoint to check answers
        e.Router.AddRoute(echo.Route {
            Method: http.MethodPost,
            Path:   "/api/submit",
            Handler: func(c echo.Context) error {
                req := c.Request()
                dao := app.Dao();

                // {{{ unmarshal request body & check answers
                var answers map[string]int;
                if err := json.NewDecoder(req.Body).Decode(&answers); err != nil { panic(err) }

                // check for user id
                userID := req.Header.Get("X-User-ID")
                attemptTime := req.Header.Get("X-Attempt-Time")
                totalTime := req.Header.Get("X-Total-Time")

                var correct bool;

                if check(uIDRegex, userID) && check(numRegex, attemptTime) && check(numRegex, totalTime) {
                    correct = checkSubmissions(dao, answers)
                } else {
                    // consider incorrect if valid custom headers missing
                    correct = false
                }
                // }}}

                // {{{ add attempt record
                tryAtoi := func(s string) int {
                    num, _ := strconv.Atoi(s);
                    if num == 0 {
                        return -1
                    } else {
                        return num
                    }
                }

                collection, err := dao.FindCollectionByNameOrId("attempts")
                if err != nil { panic(err) }

                record := models.NewRecord(collection);
                record.Set("userid", userID)
                record.Set("time", tryAtoi(attemptTime))
                record.Set("total_time", tryAtoi(totalTime))
                record.Set("correct", correct)

                if err := dao.SaveRecord(record); err != nil { panic(err) }
                // }}}

                // {{{ get keyword
                var key string
                if correct {
                    key = "successKW"
                } else {
                    key = "failureKW"
                }
                records, err := dao.FindRecordsByExpr("settings", dbx.HashExp { "key": key })
                if err != nil { panic(err) }
                // }}}

                // {{{ respond
                s := records[0].GetString("value")
                return c.JSON(http.StatusOK, map[string]any{
                    "correct": correct,
                    "keyword": s[1 : len(s) - 1],
                })
                // }}}
            },
            Middlewares: []echo.MiddlewareFunc {
                apis.ActivityLogger(app),
            },
        })
        // }}}

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
