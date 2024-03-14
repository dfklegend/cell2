package golua

import (
	"fmt"
	"os"
	"testing"
)

var gluaTests = []string{
	"base",
	"coroutine",
	"db",
	"issues",
	"os",
	"table",
	"vm",
	"math",
	"strings",
	"goto",
}

func testScriptDir(t *testing.T, tests []string, directory string) {
	if err := os.Chdir(directory); err != nil {
		t.Error(err)
	}
	defer os.Chdir("..")

	luaEngine := NewLuaEngine()
	defer luaEngine.Close()

	for _, script := range tests {
		fmt.Printf("testing %s/%s\n", directory, script)
		if err := luaEngine.DoLuaFile(fmt.Sprintf("%s.lua", script)); err != nil {
			t.Error(err)
		}
	}
}

var numActiveUserDatas int32 = 0

func TestGlua(t *testing.T) {
	testScriptDir(t, gluaTests, "_glua-tests")
}
