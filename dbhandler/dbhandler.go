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
var DATABASE_RECORD = "RECORD"
var COLLECTION_STOCK = "STOCK"
var DB, _ = NewDBHandler("/**MONGODB CONNECT URL**/")

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
	coll := db.Collection(COLLECTION_STOCK)
	return &DBHandler{client: client, db: db, coll: coll}, nil
}

func DisConnect() {
	DB.client.Disconnect(context.TODO())
}

// Stock
func (r *DBHandler) UpsertStock(stock model.NewStock) (*model.Stock, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	r.coll = r.db.Collection(COLLECTION_STOCK)
	_, err := r.coll.UpdateOne(ctx,
		bson.D{{Key: "code", Value: stock.Code}},
		bson.D{{Key: "$set", Value: stock}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return &model.Stock{
		Code: stock.Code,
		Name: stock.Name,
	}, nil
}
func (r *DBHandler) GetStockByCode(code string) ([]*model.Stock, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	r.coll = r.db.Collection(COLLECTION_STOCK)
	var res []*model.Stock
	//如果是空就返回全部項目
	var filter interface{} = bson.M{"code": code}
	if code == "" {
		filter = bson.D{{}}
	}
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
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
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	r.coll = r.db.Collection(COLLECTION_STOCK)
	res, err := r.coll.DeleteMany(ctx, bson.M{"code": code})
	if err != nil {
		return res, err
	}
	return res, nil
}

// Record
func (r *DBHandler) InsertRecord(code string, record model.DailyRecord) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	r.db = r.client.Database(DATABASE_RECORD)
	r.coll = r.db.Collection(code)
	_, err := r.coll.UpdateOne(ctx,
		bson.D{{Key: "date", Value: record.Date}},
		bson.D{{Key: "$set", Value: record}}, options.Update().SetUpsert(true))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *DBHandler) GetRecordByCode(code string, name string, date string) ([]*model.DailyRecord, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	r.db = r.client.Database(DATABASE_RECORD)
	r.coll = r.db.Collection(code)
	var res []*model.DailyRecord
	var filter interface{} = bson.M{"date": date}
	if date == "" {
		filter = bson.D{{}}
	}
	cur, err := r.coll.Find(ctx, filter)
	for cur.Next(ctx) {
		var elem *model.DailyRecord
		err := cur.Decode(&elem)
		if err != nil {
			continue
		}
		res = append(res, elem)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
