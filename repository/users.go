package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/telegram-bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepository interface {
	FindUser(userData models.UserQuery) (*models.UserData, error)
	Update(data models.UserData) (*models.UserData, error)
	Create(data models.UserData) (*models.UserData, error)
	CreateLog(data models.ActivityLog) error
	FindLog(query models.ActivityLogQuery) ([]models.ActivityLog, int64, error)
	StateCreate(data models.State) error
	GetState(userId int64) (*models.State, error)
	UpdateState(data models.State) error
	DeleteState(userId int64) error
	UpdateLog(data models.ActivityLog) error
}

type UsersRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

func (r *UsersRepositoryMongo) coll() *mongo.Collection {
	return r.ConnectionDB.Collection("employee")
}

func (r *UsersRepositoryMongo) coll2() *mongo.Collection {
	return r.ConnectionDB.Collection("activityLog")
}

func (r *UsersRepositoryMongo) coll3() *mongo.Collection {
	return r.ConnectionDB.Collection("state")
}

func (r *UsersRepositoryMongo) UpdateLog(data models.ActivityLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filter := bson.M{
		"_id":     data.ID,
		"user_id": data.UserID,
	}
	update := bson.M{
		"$set": bson.M{
			"sign_out_hour": data.SignOut,
		}, "$currentDate": bson.M{
			"_modified": true,
		},
	}

	_, err := r.coll2().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepositoryMongo) FindLog(query models.ActivityLogQuery) ([]models.ActivityLog, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	log := []models.ActivityLog{}
	filter := bson.M{}
	if query.UserID != nil {
		filter["user_id"] = query.UserID
	}
	sort := int64(-1)
	var count int64
	curr, err := r.coll2().Find(ctx, filter, &options.FindOptions{
		Sort: bson.M{"_created": sort},
	})
	if err != nil {
		return log, count, err
	}
	err = curr.All(ctx, &log)
	if err != nil {
		return log, count, err
	}
	count, err = r.coll2().CountDocuments(ctx, filter)
	return log, count, err

}

func (r *UsersRepositoryMongo) CreateLog(data models.ActivityLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data.ID = primitive.NewObjectID()
	data.Modified = time.Now()
	data.Created = data.Modified
	_, err := r.coll2().InsertOne(ctx, data)
	return err
}

func (r *UsersRepositoryMongo) FindUser(query models.UserQuery) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var user models.UserData
	filter := bson.M{}
	if query.UserID != nil {
		filter["user_id"] = *query.UserID
	}
	if query.FirstName != nil {
		filter["firstName"] = *query.FirstName
	}
	if query.LastName != nil {
		filter["lastName"] = *query.LastName
	}
	err := r.coll().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}

	return &user, err
}

func (r *UsersRepositoryMongo) Update(data models.UserData) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filter := bson.M{
		"user_id": data.UserID,
	}
	update := bson.M{
		"$set": bson.M{
			"first_name": data.FirstName,
			"last_name":  data.LastName,
			"username":   data.UserName,
		}, "$currentDate": bson.M{
			"_modified": true,
		},
	}

	_, err := r.coll().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &data, nil

}

func (r *UsersRepositoryMongo) Create(data models.UserData) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data.ID = primitive.NewObjectID()
	data.Modified = time.Now()
	data.Created = data.Modified
	_, err := r.coll().InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UsersRepositoryMongo) StateCreate(data models.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data.ID = primitive.NewObjectID()
	data.Modified = time.Now()
	data.Created = data.Modified
	_, err := r.coll3().InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepositoryMongo) GetState(userId int64) (*models.State, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var state models.State
	filter := bson.M{
		"user_id": userId,
	}
	err := r.coll3().FindOne(ctx, filter).Decode(&state)
	if err != nil {
		fmt.Println(err)
	}

	return &state, err
}

func (r *UsersRepositoryMongo) UpdateState(state models.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filter := bson.M{
		"user_id": state.UserID,
		"state":   state.State,
	}

	_, err := r.coll3().UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"sub_state": state.SubState,
			"data":      state.Data,
		}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepositoryMongo) DeleteState(userId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if _, err := r.coll3().DeleteOne(ctx, bson.M{
		"user_id": userId,
	}); err != nil {
		return err
	}

	return nil
}
