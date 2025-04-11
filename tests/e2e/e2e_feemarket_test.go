package e2e

import (
	"fmt"
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

		fmt.Println("Feemarket Params")
		fmt.Println("------")
		fmt.Println("Alpha: ", params.Params.Alpha)
		fmt.Println("Beta: ", params.Params.Beta)
		fmt.Println("Gamma: ", params.Params.Gamma)
		fmt.Println("Delta: ", params.Params.Delta)
		fmt.Println("min_base_gas_price: ", params.Params.MinBaseGasPrice)
		fmt.Println("min_learning_rate: ", params.Params.MinLearningRate)
		fmt.Println("max_learning_rate: ", params.Params.MaxLearningRate)
		fmt.Println("max_block_utilization: ", params.Params.MaxBlockUtilization)
		fmt.Println("window_size: ", params.Params.Window)
		fmt.Println("fee_denom: ", params.Params.FeeDenom)
		fmt.Println("enabled: ", params.Params.Enabled)
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
		fmt.Println("gasPrices: ", gasPrices)
		s.Require().Positive(gasPrices.Prices.AmountOf("uatone"))
		s.Require().Positive(gasPrices.Prices.AmountOf("uphoton"))
	})
}

/*
Test Gas Price change
*/
func (s *IntegrationTestSuite) testFeemarketGasPriceChange() {
	s.Run("gas price change", func() {
		fmt.Println("Called Feemarket test!!")
	})

}
