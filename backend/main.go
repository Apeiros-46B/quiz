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

func main() {
    userIDRegex := regexp.MustCompile("[A-Za-z0-9_\\-]{21}")

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
                var correct bool;
                if userIDRegex.FindString(userID) == "" {
                    // consider incorrect if no userID
                    correct = false
                } else {
                    correct = checkSubmissions(dao, answers)
                }
                // }}}

                // {{{ add attempt record
                numFromHeader := func(field string) int {
                    num, _ := strconv.Atoi(req.Header.Get(field));
                    return num
                }

                collection, err := dao.FindCollectionByNameOrId("attempts")
                if err != nil { panic(err) }

                record := models.NewRecord(collection);
                record.Set("userid", userID)
                record.Set("time", numFromHeader("X-Attempt-Time"))
                record.Set("total_time", numFromHeader("X-Total-Time"))
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
