package main

import (
	"context"
	"log"
	"time"

	"github.com/bmc-toolbox/bmclib/v2/providers/ipmitool"
	"github.com/bombsimon/logrusr/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// set BMC parameters here
	host := "10.211.132.157"
	user := "root"
	pass := "yxvZdxAQ38ZWlZ"

	l := logrus.New()
	l.Level = logrus.DebugLevel
	logger := logrusr.New(l)

	if host == "" || user == "" || pass == "" {
		log.Fatal("required host/user/pass parameters not defined")
	}

	i, err := ipmitool.New(host, user, pass, ipmitool.WithLogger(logger))
	if err != nil {
		log.Fatal("ipmi connection failed")
	}

	err = i.Open(ctx)
	if err != nil {
		log.Fatal(err, "ipmi login failed")
	}

	defer i.Close(ctx)

	output, err := i.GetIPMICiphers(ctx)
	if err != nil {
		l.Error(err)
	}
	log.Print("GetIPMICiphers returned ", output)

	output, err = i.GetSOLCiphers(ctx)
	if err != nil {
		l.Error(err)
	}
	log.Print("GetSOLCiphers returned ", output)

}
