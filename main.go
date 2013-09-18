package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
)

func setupConnection() (*xgbutil.XUtil, error) {
	X, err := xgbutil.NewConn()
	if err != nil {
		return X, err
	}
	err = xproto.ChangeWindowAttributesChecked(X.Conn(), X.RootWin(), xproto.CwEventMask, []uint32{xproto.EventMaskSubstructureRedirect | xproto.EventMaskSubstructureNotify}).Check()
	if err != nil {
		return X, err
	}
	return X, nil
}

func main() {
	fmt.Println("Welcome to Mirror!")
	X, err := setupConnection()
	if err != nil {
		fmt.Println("Mirror has encountered an error:", err)
		return
	}
	windowManger := NewWindowManager(X)
	windowManger.Run()
	return
}
