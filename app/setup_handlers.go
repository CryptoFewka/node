package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const releaseVersion = "v9.0.0"

func SetupHandlers(app *App) {
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)

		// Updated version map to the latest consensus versions from each module
		for m, mb := range app.mm.Modules {
			vm[m] = mb.ConsensusVersion()
		}
		vm[observertypes.ModuleName] = vm[observertypes.ModuleName] - 1
		vm[crosschaintypes.ModuleName] = vm[crosschaintypes.ModuleName] - 1
		SetParams(app, ctx)
		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			// Added: []string{},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// SetParams sets the default params for the observer module
// A new policy has been added for add_observer.
func SetParams(app *App, ctx sdk.Context) {
	params := app.ZetaObserverKeeper.GetParamsIsExists(ctx)
	params.AdminPolicy = observertypes.DefaultAdminPolicy()
	app.ZetaObserverKeeper.SetParams(ctx, params)
}
