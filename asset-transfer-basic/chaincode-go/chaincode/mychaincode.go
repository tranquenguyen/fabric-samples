package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract đại diện cho smart contract mới được triển khai
type SmartContract struct {
	contractapi.Contract
}

// Asset đại diện cho một tài sản trong chaincode
type Asset struct {
	ID         string  `json:"ID"`          // ID của tài sản
	Name       string  `json:"Name"`        // Tên của tài sản
	Temperature float64 `json:"Temperature"` // Nhiệt độ của tài sản
	Timestamp  string  `json:"Timestamp"`   // Thời gian ghi nhận nhiệt độ của tài sản
}

// InitLedger được gọi khi khởi tạo ledger, khởi tạo dữ liệu mẫu
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Khởi tạo một số tài sản mẫu
	assets := []Asset{
		{ID: "Sensor1", Name: "Sensor 1", Temperature: 25.5, Timestamp: "2024-04-26T10:00:00Z"},
		{ID: "Sensor2", Name: "Sensor 2", Temperature: 28.0, Timestamp: "2024-04-26T10:01:00Z"},
		{ID: "Sensor3", Name: "Sensor 3", Temperature: 24.8, Timestamp: "2024-04-26T10:02:00Z"},
	}
	// Lưu trữ các tài sản vào world state
	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// GetAllAssets trả về tất cả các tài sản hiện có trong world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}

// CreateAsset tạo một tài sản mới trong world state
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, name string) error {
	// Kiểm tra xem tài sản đã tồn tại chưa
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	// Nếu tài sản đã tồn tại, trả về lỗi
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}
	// Tạo một tài sản mới
	asset := Asset{
		ID:   id,
		Name: name,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	// Lưu trữ tài sản vào world state
	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset đọc thông tin của một tài sản từ world state
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}
	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// AssetExists kiểm tra xem một tài sản có tồn tại trong world state không
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

// DeleteAsset xóa một tài sản khỏi world state
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	// Kiểm tra xem tài sản có tồn tại không
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	// Nếu tài sản không tồn tại, trả về lỗi
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	// Xóa tài sản từ world state
	return ctx.GetStub().DelState(id)
}

// Xóa tất cả các tài sản từ world state
func (s *SmartContract) DeleteAllAssets(ctx contractapi.TransactionContextInterface) error {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return err
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}
		err = ctx.GetStub().DelState(queryResponse.Key)
		if err != nil {
			return err
		}
	}
	return nil
}


