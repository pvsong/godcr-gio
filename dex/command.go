package dex

import (
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex/encode"
)

func (d *Dex) InitializeClient(apppasswd string) {
	err := d.core.InitializeClient([]byte(apppasswd))
	if err != nil {
		fmt.Printf("InitializeClient error: %v", err)
	}
}

func (d *Dex) SupportedAsset() map[uint32]*core.SupportedAsset {
	return d.core.SupportedAssets()
}

func (d *Dex) IsInitialized() bool {
	if _, err := d.core.IsInitialized(); err != nil {
		return false
	}

	return true
}

func (d *Dex) GetUser() *core.User {
	u := d.core.User()
	return u
}

func (d *Dex) GetDefaultWalletConfig() map[string]string {
	cfg, err := d.core.AutoWalletConfig(42)
	if err != nil {
		return nil
	}
	return cfg
}

func (d *Dex) GetDEXConfig(addr string, cert string) {
	go func() {
		cx, err := d.core.GetDEXConfig(addr, []byte(cert))
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("[SUCCESS:GetDEXConfig]", cx)
		return
	}()
}

func (d *Dex) Register(apppasswd string) {
	go func() {
		form := &core.RegisterForm{
			AppPass: []byte(apppasswd),
			Addr:    "string",
		}
		result, err := d.core.Register(form)
		log.Info(result)

		if err != nil {
			fmt.Printf("InitializeClient error: %v", err)
		}
	}()
}

type newWalletForm struct {
	AssetID uint32
	Config  map[string]string
	Pass    encode.PassBytes
	AppPW   encode.PassBytes
}

func (d *Dex) Login(apppasswd string) {
	go func() {
		_, err := d.core.Login([]byte(apppasswd))
		if err != nil {
			fmt.Printf("InitializeClient error: %v", err)
		}
	}()
}

func (d *Dex) AddNewWallet(form newWalletForm) {
	go func() {
		has := d.core.WalletState(form.AssetID) != nil
		if has {
			return
		}

		// Wallet does not exist yet. Try to create it.
		err := d.core.CreateWallet(form.AppPW, form.Pass, &core.WalletForm{
			AssetID: form.AssetID,
			Config:  form.Config,
		})
		if err != nil {
			return
		}

		// Send response ok
	}()
}
