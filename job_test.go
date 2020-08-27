package crontab_test

import (
	"sync"
	"testing"
	"time"

	"github.com/shved/crontab"
)

func TestJobError(t *testing.T) {

	ctab := crontab.New()
	ctab.Start()

	if err := ctab.AddJob("* * * * *", "asdf1", myFunc, 10); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", "asdf2", nil); err == nil {
		t.Error("This AddJob should return Error, fn is nil")
	}

	var x int
	if err := ctab.AddJob("* * * * *", "asdf3", x); err == nil {
		t.Error("This AddJob should return Error, fn is not func kind")
	}

	if err := ctab.AddJob("* * * * *", "asdf4", myFunc2, "s", 10, 12); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", "asdf5", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, args are not the correct type")
	}

	if err := ctab.AddJob("* * * * * *", "asdf6", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, syntax error")
	}

	// custom types and interfaces as function params
	var m MyTypeInterface
	if err := ctab.AddJob("* * * * *", "asdf7", myFuncStruct, m); err != nil {
		t.Error(err)
	}

	if err := ctab.AddJob("* * * * *", "asdf8", myFuncInterface, m); err != nil {
		t.Error(err)
	}

	var mwo MyTypeNoInterface
	if err := ctab.AddJob("* * * * *", "asdf9", myFuncInterface, mwo); err == nil {
		t.Error("This should return error, type that don't implements interface assigned as param")
	}

	if err := ctab.AddJob("* * * * *", "asdf9", nil); err == nil {
		t.Error("This should return error, name is already registered")
	}

	ctab.Shutdown()
}

var testN int
var testS string

func TestCrontab(t *testing.T) {
	testN = 0
	testS = ""

	ctab := crontab.Fake(2) // fake crontab wiht 2sec timer to speed up test
	ctab.Start()

	var wg sync.WaitGroup
	wg.Add(2)

	if err := ctab.AddJob("* * * * *", "asdf1", func() { testN++; wg.Done() }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", "asdf2", func(s string) { testS = s; wg.Done() }, "param"); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}

	if testN != 1 {
		t.Error("func 1 not executed as scheduled")
	}

	if testS != "param" {
		t.Error("func 2 not executed as scheduled")
	}
	ctab.Shutdown()
}

func TestRun(t *testing.T) {
	testN = 0
	testS = "test"

	ctab := crontab.New()
	ctab.Start()

	if err := ctab.AddJob("* * * * *", "asdf1", func() { testN++ }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", "asdf2", func(s string) { testS = s }, "param"); err != nil {
		t.Fatal(err)
	}

	if err := ctab.Run("asdf1"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if testN != 1 {
		t.Error("func not executed on Run()")
	}

	if testS != "test" {
		t.Error("wrong func executed on Run()")
	}

	if err := ctab.Run("missing_job"); err == nil {
		t.Error("invoking missing func name doesnt throw an error on Run()")
	}
}

func TestRunAll(t *testing.T) {
	testN = 0
	testS = ""

	ctab := crontab.New()
	ctab.Start()

	if err := ctab.AddJob("* * * * *", "asdf1", func() { testN++ }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", "asdf2", func(s string) { testS = s }, "param"); err != nil {
		t.Fatal(err)
	}

	ctab.RunAll()
	time.Sleep(time.Second)

	if testN != 1 {
		t.Error("func not executed on RunAll()")
	}

	if testS != "param" {
		t.Error("func not executed on RunAll() or arg not passed")
	}

	ctab.Clear()
	ctab.RunAll()

	if testN != 1 {
		t.Error("Jobs not cleared")
	}

	if testS != "param" {
		t.Error("Jobs not cleared")
	}

	ctab.Shutdown()
}
