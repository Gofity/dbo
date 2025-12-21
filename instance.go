package dbo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SeederEntry interface {
	Name() string
	Handler(db *gorm.DB) (err error)
}

type Instance struct {
	opts Options
	db   *gorm.DB
}

func (x *Instance) Get(options ...O) (db *gorm.DB, err error) {
	opts := x.dbOptions(options...)

	if db = opts.DB; db == nil {
		db = NewSession(x.db)
	}

	if db == nil {
		err = errors.New("Database connection not found")
		return
	}

	// Add Clauses
	if len(opts.Clauses) > 0 {
		db = db.Clauses(opts.Clauses...).Session(&gorm.Session{})
	}

	// Add Scopes
	if len(opts.Scopes) > 0 {
		db = db.Scopes(opts.Scopes...).Session(&gorm.Session{})
	}

	return
}

func (x *Instance) WithCancel(options ...O) (db *gorm.DB, cancel context.CancelFunc, err error) {
	opts := x.dbOptions(options...)

	db, err = x.Get(opts)

	if err != nil {
		return
	}

	ctx, cancel := opts.newContext()

	db = db.WithContext(ctx)

	return
}

func (x *Instance) Migrate(models ...any) (err error) {
	fmt.Print("Running Database Migration... ")

	db, err := x.Get()

	if err != nil {
		return
	}

	switch x.opts.getDriver() {
	case DRIVER_MYSQL:
		db = db.Set("gorm:table_options", fmt.Sprintf(
			"ENGINE=%s CHARSET=%s COLLATE=%s",
			x.opts.getEngine(),
			x.opts.getCharset(),
			x.opts.getCollation(),
		))
	}

	err = db.AutoMigrate(models...)

	if err != nil {
		return
	}

	fmt.Println("Completed!")
	fmt.Println()

	return
}

func (x *Instance) Seed(entries ...SeederEntry) (err error) {
	fmt.Println("Running Database Seeders...")

	db, err := x.Get()

	if err != nil {
		return
	}

	for _, entry := range entries {
		err = x.seed(NewSession(db), entry)

		if err != nil {
			return
		}
	}

	fmt.Println("Database Seeders Completed!")
	fmt.Println()

	return
}

func (x *Instance) seed(db *gorm.DB, entry SeederEntry) (err error) {
	name := strings.TrimSpace(entry.Name())

	fmt.Printf("Seeding %v... ", name)

	err = entry.Handler(db)

	if err != nil {
		return
	}

	fmt.Println("Completed!")

	return
}

func (x *Instance) dbOptions(options ...O) O {
	if len(options) > 0 {
		return options[0]
	}

	return O{}
}
