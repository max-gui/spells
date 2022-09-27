package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/max-gui/spells/internal/iac/altconfig"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/pkg/constset"
)

var plaintext, cryptedHexText, md5Hex string
var router *gin.Engine

func setup() {
	gin.SetMode(gin.TestMode)
	plaintext = "123"
	cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	md5Hex = "202cb962ac59075b964b07152d234b70"
	constset.StartupInit(nil, context.Background())
	// abstestpath = confgen.Makeconfiglist()

	router = SetupRouter()
	// fmt.Println(config.AppSetting.JwtSecret)
	// fmt.Println("Before all tests")
}

func teardown() {

}

func Test_ArchDef_commit_check(t *testing.T) {
	// router := SetupRouter()

	// {
	// 	"commiter":"aaa",
	// 	"branch":"ddd",
	// 	"changedFiles": [
	// 		{"path":"","status":"A","content":""}
	// 	],
	// 	"commitIds":[""],
	// 	"messages":[""],
	// 	"repositoryName":"name",
	// 	"namespace":"组名"
	// }
	jsond := altconfig.CommitCheckHookInfo{}
	jsond.Branch = "branch"
	jsond.CommitIds = []string{"0", "1"}
	jsond.Commiter = "tester"
	jsond.Messages = []string{"m0", "m1"}
	jsond.Namespace = "ns"
	jsond.RepositoryName = "repo"
	jsond.Type = "git"
	jsond.ChangedFiles = []archfig.FileContentInfo{}

	jsonmap := make(map[string]interface{})
	jsonmap["data"] = plaintext
	jsonByte, _ := json.Marshal(jsonmap)
	// body, writer := body4PostFile("yamls"+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/encrypt2Hex", bytes.NewReader(jsonByte))
	// req.Header.Add("Content-type", writer.FormDataContentType())

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	resbody, _ := io.ReadAll(result.Body)

	resstr := string(resbody)
	resjsonmap := make(map[string]interface{})
	json.Unmarshal(resbody, &resjsonmap)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, cryptedHexText, resjsonmap["data"])

	t.Logf("Test_DecryptHexonline result is:\n%s", resstr)
}

func TestMain(m *testing.M) {
	setup()
	// constset.StartupInit()
	// sendconfig2consul()
	// configgen.Getconfig = getTestConfig

	exitCode := m.Run()
	teardown()
	// // 退出
	os.Exit(exitCode)
}
