package entities

type Investor struct {
	ID             string
	Name           string
	AssetPositions []*InvestorAssetPosition
}

type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:             id,
		AssetPositions: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPositions = append(i.AssetPositions, assetPosition)
}

func (i *Investor) UpdateAssetPosition(assetID string, shares int) {
	assetPosition := i.GetAssetPosition(assetID)

	if assetPosition == nil {
		i.AssetPositions = append(i.AssetPositions, NewAssetPosition(assetID, shares))
	} else {
		assetPosition.Shares += shares
	}
}

func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPositions {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}
	return nil
}

func NewAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  shares,
	}
}
