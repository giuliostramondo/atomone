package e2e

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"time"
)

/*
Test Feemarket Queries:
- params
- state
- gas_price/{denom}
- gas_prices
*/
func (s *IntegrationTestSuite) testFeemarketQuery() {
	s.Run("feemarket test params", func() {
		var (
			c             = s.chainA
			valIdx        = 0
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)
		params := s.queryFeemarketParams(chainEndpoint)
		var maxBlockUtilization uint64 = 30000000
		var window_size uint64 = 1
		s.Require().Equal("0.000000000000000000", params.Params.Alpha.String())
		s.Require().Equal("1.000000000000000000", params.Params.Beta.String())
		s.Require().Equal("0.000000000000000000", params.Params.Gamma.String())
		s.Require().Equal("0.000000000000000000", params.Params.Delta.String())
		s.Require().Equal("0.000010000000000000", params.Params.MinBaseGasPrice.String())
		s.Require().Equal("0.125000000000000000", params.Params.MinLearningRate.String())
		s.Require().Equal("0.125000000000000000", params.Params.MaxLearningRate.String())
		s.Require().Equal("0.125000000000000000", params.Params.MaxLearningRate.String())
		s.Require().Equal(maxBlockUtilization, params.Params.MaxBlockUtilization)
		s.Require().Equal(window_size, params.Params.Window)
		s.Require().Equal("uphoton", params.Params.FeeDenom)
		s.Require().True(params.Params.Enabled)

		s.T().Logf("Feemarket Params")
		s.T().Logf("------")
		s.T().Logf("Alpha: %s", params.Params.Alpha)
		s.T().Logf("Beta: %s", params.Params.Beta)
		s.T().Logf("Gamma: %s", params.Params.Gamma)
		s.T().Logf("Delta: %s", params.Params.Delta)
		s.T().Logf("min_base_gas_price: %s", params.Params.MinBaseGasPrice)
		s.T().Logf("min_learning_rate: %s", params.Params.MinLearningRate)
		s.T().Logf("max_learning_rate: %s", params.Params.MaxLearningRate)
		s.T().Logf("max_block_utilization: %d", params.Params.MaxBlockUtilization)
		s.T().Logf("window_size: %d", params.Params.Window)
		s.T().Logf("fee_denom: %s", params.Params.FeeDenom)
		s.T().Logf("enabled: %t", params.Params.Enabled)
	})
	s.Run("feemarket test state", func() {
		var (
			c             = s.chainA
			valIdx        = 0
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)
		params := s.queryFeemarketState(chainEndpoint)
		fmt.Println("State")
		fmt.Println("-----")
		fmt.Println("BaseGasPrice: ", params.State.BaseGasPrice)
		fmt.Println("Window: ", params.State.Window)
		fmt.Println("LearningRate: ", params.State.LearningRate)
		fmt.Println("Index: ", params.State.Index)
	})

	s.Run("feemarket test get price", func() {
		var (
			c             = s.chainA
			valIdx        = 0
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)
		atoneGasPrice := s.queryFeemarketGasPrice(chainEndpoint, "uatone")
		s.Require().Equal("uatone", atoneGasPrice.Price.Denom)
		photonGasPrice := s.queryFeemarketGasPrice(chainEndpoint, "uphoton")
		s.Require().Equal("uphoton", photonGasPrice.Price.Denom)
		fmt.Println("atoneGasPrice: ", atoneGasPrice)
		fmt.Println("photonGasPrice: ", photonGasPrice)
	})

	s.Run("feemarket test get prices", func() {
		var (
			c             = s.chainA
			valIdx        = 0
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)
		gasPrices := s.queryFeemarketGasPrices(chainEndpoint)
		s.T().Logf("gasPrices: %s", gasPrices)
		atoneAmount := gasPrices.Prices.AmountOf("uatone")
		photonAmount := gasPrices.Prices.AmountOf("uphoton")
		s.Require().True(atoneAmount.IsPositive())
		s.Require().True(photonAmount.IsPositive())
	})
}

/*
Test Gas Price change
*/
func (s *IntegrationTestSuite) testFeemarketGasPriceChange() {
	s.Run("gas price change", func() {
		var (
			c             = s.chainA
			valIdx        = 0
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)
		gasPricesInitial := s.queryFeemarketGasPrices(chainEndpoint)
		s.T().Logf("Initial gasPrices: %s", gasPricesInitial)
		// define one sender and two recipient accounts
		sender, _ := c.genesisAccounts[0].keyInfo.GetAddress()

		var beforeAccountBalances,
			afterAccountBalances []sdk.Coin

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				for i := range len(c.genesisAccounts) {
					accountID := i
					address, _ := c.genesisAccounts[accountID].keyInfo.GetAddress()
					addressUAtoneBalance, err := getSpecificBalance(chainEndpoint, address.String(), uatoneDenom)
					beforeAccountBalances = append(beforeAccountBalances, addressUAtoneBalance)
					s.Require().NoError(err)
				}

				balanceValid := beforeAccountBalances[0].IsValid()

				for i := range len(c.genesisAccounts) {
					balanceValid = balanceValid && beforeAccountBalances[i].IsValid()
				}
				return balanceValid
			},
			10*time.Second,
			time.Second,
		)

		s.T().Logf("Total Number Of account: %d", len(c.genesisAccounts))

		s.T().Logf("Initial Account Balances")
		for i := range len(c.genesisAccounts) {
			accountID := i
			s.T().Logf("Account %d: %d", i, beforeAccountBalances[accountID].Amount)
		}

		var destAccounts []string

		//tokenAmount = sdk.NewInt64Coin(uatoneDenom, 100_000) // 0.1atone
		txNumber := 2
		for i := range txNumber {
			accountID := i%len(c.genesisAccounts) + 1
			address, _ := c.genesisAccounts[accountID].keyInfo.GetAddress()
			destAccounts = append(destAccounts, address.String())
		}
		// alice sends tokens to bob and charlie, at once
		resp := s.execBankMultiSendAndReturn(s.chainA, valIdx, sender.String(),
			destAccounts, tokenAmount.String(), true)

		s.T().Logf("Response :\n %s", resp.String())

		getRequest := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s",
			chainEndpoint, resp.TxHash)

		s.T().Logf("HTTP get request :\n %s", getRequest)
		body, errHttp := httpGet(getRequest)
		var txResp tx.GetTxResponse
		s.T().Logf("Tx Err Response:\n%s", errHttp)
		s.T().Logf("Tx RAW Response:\n%s", body)
		if err := cdc.UnmarshalJSON(body, &txResp); err != nil {
			s.T().Logf("failed to read response body: %s", err)
		}
		s.T().Logf("Tx Response:\n%s", txResp)
		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				for i := range len(c.genesisAccounts) {
					accountID := i
					address, _ := c.genesisAccounts[accountID].keyInfo.GetAddress()
					addressUAtoneBalance, err := getSpecificBalance(chainEndpoint, address.String(), uatoneDenom)
					afterAccountBalances = append(beforeAccountBalances, addressUAtoneBalance)
					s.Require().NoError(err)
				}

				balanceValid := afterAccountBalances[0].IsValid()
				for i := range len(c.genesisAccounts) {
					balanceValid = balanceValid && afterAccountBalances[i].IsValid()
				}

				balancesAreDifferent := false
				for i := range len(c.genesisAccounts) {
					balancesAreDifferent = balancesAreDifferent ||
						!afterAccountBalances[i].Amount.Equal(beforeAccountBalances[i].Amount)
				}

				return balanceValid //&& balancesAreDifferent
			},
			10*time.Second,
			time.Second,
		)

		s.T().Logf("Final Account Balances")
		for i := range len(c.genesisAccounts) {
			accountID := i
			s.T().Logf("Account %d: %d", i, afterAccountBalances[accountID].Amount)
		}

		gasPricesFinal := s.queryFeemarketGasPrices(chainEndpoint)
		s.T().Logf("Final gasPrices: %s", gasPricesFinal)
	})

}
