// package config

// import (
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/cloudwego/hertz/pkg/common/test/assert"
// )

// func TestInitConfig(t *testing.T) {
// 	// Set environment variable for testing
// 	os.Setenv("ENV", "prod")
// 	InitConfig()
// 	assert.DeepEqual(t, specificCfg, 
// 		&SpecificConfig{
// 			Env:    "prod",
// 			DB_DSN: "postgres://prod_user:pass@prod-db:5432/proddb",
// 			Domain: "skylar27.com",
// 		},
// 	)
// 	fmt.Println(generalCfg)
// 	assert.DeepEqual(t, generalCfg,
// 		&GeneralConfig{
// 			JWT_Secret_Key: "secret-secret-secret",
// 		},
// 	)
// }
package config