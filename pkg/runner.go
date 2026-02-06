// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"os"
	"runtime"
	"syscall"
	"os/signal"

	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/app/proxyman/inbound"
)

// Runs xray-core instance (non-blocking)
func (v2 *V2utils) Run_Xray() error {
	var err error
	var cf *core.Config

	cf, err = v2.CFG.Build()
	if nil != err {
		return err
	}

	if v2.Xray_instance, err = core.New(cf); nil != err {
		return err
	}
	// Cleanup sh** we have done so far to make the config
	runtime.GC()

	if err = v2.Xray_instance.Start(); nil != err {
		return err
	}
	return nil
}

// Run_Xray (blocking)
func (v2 *V2utils) Exec_Xray() error {
	if e := v2.Run_Xray(); nil != e {
		return e
	}
	defer v2.Kill_Xray();

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
	return nil
}

// Do NOT use v2.Xray_instance after this call
func (v2 V2utils) Kill_Xray() {
	if nil != v2.Xray_instance {
		if v2.Xray_instance.IsRunning() {
			v2.Xray_instance.Close()
			v2.Xray_instance = nil
		}
	}
}
