package main
import (
	"os"
	"fmt"
)
func InitEnviromentVars()  {
	env := os.Getenv("LOGNAME")
	fmt.Println(env)
}