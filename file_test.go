package logs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	BeegoLogs "github.com/astaxie/beego/logs"
)

func TestFile1(t *testing.T) {
	log := BeegoLogs.NewLogger(10000)
	err := log.SetLogger("multi_file", `{}`)
	if err == nil {
		t.Fatal("shoud have err")
	}
}

func TestFile2(t *testing.T) {
	debugFileName := "test_file2.debug.log"
	os.Remove(debugFileName)
	defer os.Remove(debugFileName)

	log := BeegoLogs.NewLogger(10000)
	err := log.SetLogger("multi_file", fmt.Sprintf(
		`{
			"levelfiles": [{
				"levelname": "debug",
				"filename": "%s"
			}]
		}`, debugFileName))
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("debug msg")
	b, err := hasContent(debugFileName)
	if !b || err != nil {
		t.Fatalf("log file[%s] doesn't exist", debugFileName)
	}
}

func TestFile3(t *testing.T) {
	debugFileName := "test_file3.debug.log"
	infoFileName := "test_file3.info.log"
	os.Remove(debugFileName)
	os.Remove(infoFileName)
	defer os.Remove(debugFileName)
	defer os.Remove(infoFileName)

	log := BeegoLogs.NewLogger(10000)
	err := log.SetLogger("multi_file", fmt.Sprintf(
		`{
			"levelfiles": [{
				"levelname": "debug",
				"filename": "%s"
			},{
				"levelname": "info",
				"filename": "%s"
			}]
		}`, debugFileName, infoFileName))
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("debug msg")
	{
		b, err := hasContent(debugFileName)
		if !b || err != nil {
			t.Fatalf("log file[%s] doesn't exist", debugFileName)
		}
	}

	log.Info("info msg")
	{
		b, err := hasContent(infoFileName)
		if !b || err != nil {
			t.Fatalf("log file[%s] doesn't exist", infoFileName)
		}
	}
}

func TestFile4(t *testing.T) {
	debugFileName := "test_file4.debug.log"
	infoFileName := "test_file4.info.log"
	os.Remove(debugFileName)
	os.Remove(infoFileName)
	defer os.Remove(debugFileName)
	defer os.Remove(infoFileName)

	log := BeegoLogs.NewLogger(10000)
	err := log.SetLogger("multi_file", fmt.Sprintf(
		`{
			"levelname": "info",
			"levelfiles": [{
				"levelname": "debug",
				"filename": "%s"
			},{
				"levelname": "info",
				"filename": "%s"
			}]
		}`, debugFileName, infoFileName))
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("debug msg")
	{
		b, err := hasContent(debugFileName)
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.Fatalf("log file[%s] should not have content", debugFileName)
		}
	}

	log.Info("info msg")
	{
		b, err := hasContent(infoFileName)
		if err != nil {
			t.Fatal(err)
		}
		if !b {
			t.Fatalf("log file[%s] should have content", infoFileName)
		}
	}
}

func hasContent(path string) (bool, error) {
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	return len(contentBytes) > 0, nil
}
