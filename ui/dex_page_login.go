package ui

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/dex"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageDex = "Dex"

type DexPage struct {
	dexc           *dex.Dex
	theme          *decredmaterial.Theme
	passwordEditor decredmaterial.Editor
	supportedAsset []*core.SupportedAsset
	user           *core.User

	createPassword decredmaterial.Button
	login          decredmaterial.Button
	addWallConfig  decredmaterial.Button
	addDexServer   decredmaterial.Button

	isAppInitialized bool
}

func (win *Window) DexPage(common pageCommon) layout.Widget {
	pg := &DexPage{
		dexc:  win.dexc,
		theme: common.theme,
		user:  new(core.User),

		createPassword: common.theme.Button(new(widget.Clickable), "Create password"),
		login:          common.theme.Button(new(widget.Clickable), "Login"),
		addWallConfig:  common.theme.Button(new(widget.Clickable), "Add wallet config"),
		addDexServer:   common.theme.Button(new(widget.Clickable), "Add a DEX"),
	}
	pg.passwordEditor = win.theme.Editor(new(widget.Editor), "App password")
	pg.passwordEditor.Editor.SetText("")
	pg.passwordEditor.Editor.SingleLine = true

	// Get initial values
	pg.isAppInitialized = pg.dexc.IsInitialized()
	for _, v := range pg.dexc.SupportedAsset() {
		pg.supportedAsset = append(pg.supportedAsset, v)
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *DexPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

			// Init app
			layout.Rigid(func(gtx C) D {
				if pg.isAppInitialized {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.passwordEditor.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return pg.createPassword.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.theme.Separator().Layout(gtx)
			}),

			// Login
			layout.Rigid(func(gtx C) D {
				if pg.user.Initialized {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.passwordEditor.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return pg.login.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.theme.Separator().Layout(gtx)
			}),

			// Register
			layout.Rigid(func(gtx C) D {
				return pg.addWallConfig.Layout(gtx)
			}),
		)
	}

	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *DexPage) handle(common pageCommon) {

	if pg.createPassword.Button.Clicked() {
		pg.dexc.InitializeClient(pg.passwordEditor.Editor.Text())
		pg.user = pg.dexc.GetUser()
	}

	if pg.login.Button.Clicked() {
		pg.dexc.Login(pg.passwordEditor.Editor.Text())
		u := pg.dexc.GetUser()
		pg.user = u
	}

	if pg.login.Button.Clicked() {
		pg.user = pg.dexc.GetUser()
	}

	if pg.addWallConfig.Button.Clicked() {
		df := pg.dexc.GetDefaultWalletConfig()
		log.Info("df.Password", df["password"])
		log.Info("df.RpcCert", df["rpccert"])
		log.Info("df.RpcKey", df["rpckey"])
		log.Info("df.RpcListen", df["rpclisten"])
		log.Info("df.Username", df["username"])
		pg.dexc.GetDEXConfig("http://127.0.0.1:7232", "Your cert...")
	}

}
