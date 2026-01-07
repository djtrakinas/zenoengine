package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"zeno/internal/app"
	"zeno/pkg/dbmanager"
	"zeno/pkg/engine"
	"zeno/pkg/logger"
	"zeno/pkg/worker"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func HandleRun(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: zeno run <path/to/script.zl>")
		os.Exit(1)
	}
	godotenv.Load()
	logger.Setup("development")
	path := args[0]
	root, err := engine.LoadScript(path)
	if err != nil {
		fmt.Printf("❌ Syntax Error: %v\n", err)
		os.Exit(1)
	}

	dbMgr := dbmanager.NewDBManager()
	eng := engine.NewEngine()

	// Setup DB Connection
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "mysql"
	}

	var dsn string
	if dbDriver == "sqlite" {
		dsn = os.Getenv("DB_NAME")
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	}

	if err := dbMgr.AddConnection("default", dbDriver, dsn, 10, 5); err != nil {
		fmt.Printf("❌ Fatal: DB Connection Failed: %v\n", err)
		os.Exit(1)
	}

	// Auto-detect additional DBs (Copied from zeno.go)
	envVars := os.Environ()
	detectedDBs := make(map[string]bool)
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		if strings.HasPrefix(key, "DB_") && strings.HasSuffix(key, "_HOST") {
			if key == "DB_HOST" {
				continue
			}
			dbName := strings.ToLower(strings.TrimSuffix(strings.TrimPrefix(key, "DB_"), "_HOST"))
			if dbName != "" {
				detectedDBs[dbName] = true
			}
		}
	}
	for dbName := range detectedDBs {
		prefix := "DB_" + strings.ToUpper(dbName) + "_"
		// Additional DBs are assumed to be MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv(prefix+"USER"), os.Getenv(prefix+"PASS"), os.Getenv(prefix+"HOST"), os.Getenv(prefix+"NAME"))
		if err := dbMgr.AddConnection(dbName, "mysql", dsn, 10, 5); err != nil {
			fmt.Printf("⚠️  Failed to connect to database %s: %v\n", dbName, err)
		} else {
			fmt.Printf("✅ Additional Database Connected! db=%s\n", dbName)
		}
	}

	// Gunakan helper registry yang baru dibuat
	queue := worker.NewDBQueue(dbMgr, "default")
	r := chi.NewRouter()
	app.RegisterAllSlots(eng, r, dbMgr, queue, nil)

	if err := eng.Execute(context.Background(), root, engine.NewScope(nil)); err != nil {
		fmt.Printf("❌ Execution Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
