package entity

type Investor struct {
	ID	 string
	Name string
	AssetPositions []*InvestorAssetPosition
}

func NewInvestor (id string) *Investor {
	return &Investor{
		ID: id,
		AssetPositions: []*InvestorAssetPosition{},
	}
}

func (investor *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	investor.AssetPositions = append(investor.AssetPositions, assetPosition)
}

func (investor *Investor) UpdateAssetPosition(assetId string, qtnShares int) {
	assetPostion := investor.GetAssetPosition(assetId)
	if assetPostion == nil {
		investor.AssetPositions = append(investor.AssetPositions, NewInvestorAssetPosition(assetId, qtnShares))
	} else {
		assetPostion.Shares += qtnShares
	}
}

func (investor *Investor) GetAssetPosition(assetId string) *InvestorAssetPosition {
	for _, assetPosition := range investor.AssetPositions {
		if assetPosition.AssetID == assetId {
			return assetPosition
		}
	}

	return nil
}


type InvestorAssetPosition struct {
	AssetID string // ID da ação
	Shares  int    // Quantidade de cotas dessa ação
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares: shares,
	}
}