package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"os"
	"path"
)

var Casbin *casbin.Enforcer

func Dial(db *gorm.DB) error {
	dir, _ := os.Getwd()
	a, _ := gormadapter.NewAdapterByDB(db)
	Casbin, _ = casbin.NewEnforcer(path.Join(dir, "/internal/casbin/model.conf"), a)

	// Load the policy from DB.
	return Casbin.LoadPolicy()
}