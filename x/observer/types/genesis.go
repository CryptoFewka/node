package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultGenesis returns the default observer genesis state
func DefaultGenesis() *GenesisState {
	params := DefaultParams()
	return &GenesisState{
		Params:            &params,
		Ballots:           nil,
		Observers:         nil,
		NodeAccountList:   []*NodeAccount{},
		PermissionFlags:   &PermissionFlags{IsInboundEnabled: true, IsOutboundEnabled: true},
		Keygen:            nil,
		LastObserverCount: nil,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if gs.Params != nil {
		err := gs.Params.Validate()
		if err != nil {
			return err
		}
	}
	// Check for duplicated index in nodeAccount
	nodeAccountIndexMap := make(map[string]bool)

	for _, elem := range gs.NodeAccountList {
		if _, ok := nodeAccountIndexMap[elem.GetOperator()]; ok {
			return fmt.Errorf("duplicated index for nodeAccount")
		}
		nodeAccountIndexMap[elem.GetOperator()] = true
	}

	if err := VerifyObserverMapper(gs.Observers); err != nil {
		return err
	}

	return nil
}

func GetGenesisStateFromAppState(marshaler codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := marshaler.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}
