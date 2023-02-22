package dbhandler

import (
	"context"
	"time"

	"github.com/kwangsing3/stock-bar/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DATABASE = "DEBUG"
var COLLECTION = "STOCK"
var DB, _ = NewDBHandler("mongodb+srv://genesis:xRncyo2dVDcvPiJQ@cluster0.jm9ahx2.mongodb.net/test")

type DBHandler struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

func NewDBHandler(srv string) (*DBHandler, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(srv))
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}
	db := client.Database(DATABASE)
	coll := db.Collection(COLLECTION)
	return &DBHandler{client: client, db: db, coll: coll}, nil
}

func DisConnect() {
	DB.client.Disconnect(context.TODO())
}

// Stock
func (r *DBHandler) UpsertStock(stock model.NewStock) (*model.Stock, error) {
	res := model.Stock{
		Name:             stock.Name,
		Code:             stock.Code,
		HistoricalRecord: []*model.DailyRecord{},
	}
	_, err := r.coll.UpdateOne(context.TODO(),
		bson.D{{Key: "code", Value: stock.Code}},
		bson.D{{Key: "$set", Value: res}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (r *DBHandler) GetStockByCode(code string) ([]*model.Stock, error) {
	var res []*model.Stock
	//如果是空就返回全部項目
	var filter interface{} = bson.M{"code": code}
	if code == "" {
		filter = bson.D{{}}
	}
	cur, err := r.coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem *model.Stock
		err := cur.Decode(&elem)
		if err != nil {
			continue
		}
		res = append(res, elem)
	}
	return res, nil
}
func (r *DBHandler) DeleteStock(code string) (*mongo.DeleteResult, error) {
	res, err := r.coll.DeleteMany(context.TODO(), bson.M{"code": code})
	if err != nil {
		return res, err
	}
	return res, nil
}

// Record
func (r *DBHandler) InsertRecord(code string, record model.DailyRecord) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"code": code}
	arrayFilters := options.ArrayFilters{Filters: bson.A{bson.M{"x.date": record.Date}}}
	upsert := true
	opts := options.UpdateOptions{
		ArrayFilters: &arrayFilters,
		Upsert:       &upsert,
	}
	update := bson.M{
		"$set": bson.M{
			"historicalrecord.$[x]": record,
		},
	}
	_, err := r.coll.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *DBHandler) UpdateRecord(code string, record model.NewRecord) error {
	filter := bson.M{"code": code}
	update := bson.M{
		"$set": bson.M{
			"historicalrecord": record,
		},
	}
	_, err := r.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBHandler) GetRecordByCode(code string, name string, date string) ([]*model.DailyRecord, error) {
	var res []*model.DailyRecord
	filter := bson.M{"code": code}
	var dd *model.Stock
	err := r.coll.FindOne(context.TODO(), filter).Decode(&dd)

	lens := len(dd.HistoricalRecord)
	for i := 0; i < lens; i++ {
		if dd.HistoricalRecord[i].Date == date {
			res = append(res, dd.HistoricalRecord[i])
		}
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
