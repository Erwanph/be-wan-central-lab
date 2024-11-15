package repository

import (
	"context"
	"time"

	"github.com/Erwanph/be-wan-central-lab/internal/entity"
	"github.com/Erwanph/be-wan-central-lab/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DB *mongo.Client
}

func NewUserRepository(db *mongo.Client) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateCustomUser(ctx context.Context, user *entity.User) error {
	collection := r.DB.Database("digital-voter").Collection("users")
	user.CreatedAt = time.Now()
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateDefaultUser(ctx context.Context, user *entity.User) error {
	//default user
	collection := r.DB.Database("digital-voter").Collection("users")
	user.CreatedAt = time.Now()
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CountByEmail(ctx context.Context, username string) (int, error) {
	collection := r.DB.Database("digital-voter").Collection("users")
	count, err := collection.CountDocuments(ctx, map[string]interface{}{"email": username})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	collection := r.DB.Database("digital-voter").Collection("users")

	err := collection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) VerifiedEmailUser(ctx context.Context, user *entity.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"is_email_verified": user.IsEmailVerified,
			"otp":               "",
			"otp_expires_at":    nil,
			"updated_at":        util.NowInWIB(),
		},
	}
	_, err := r.DB.Database("digital-voter").Collection("users").UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepository) FindUserByResetToken(ctx context.Context, token string) (*entity.User, error) {
	collection := r.DB.Database("digital-voter").Collection("users")
	filter := bson.M{"reset_token": token}
	var user entity.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) FindByID(ctx context.Context, _id primitive.ObjectID) (*entity.User, error) {
	user := &entity.User{}
	collection := r.DB.Database("digital-voter").Collection("users")
	err := collection.FindOne(ctx, map[string]interface{}{"_id": _id}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateResetToken(ctx context.Context, userID primitive.ObjectID, token string, expiry time.Time) error {
	collection := r.DB.Database("digital-voter").Collection("users")

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"reset_token":        token,
			"reset_token_expiry": expiry,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *UserRepository) UpdatePassword(ctx context.Context, user *entity.User) error {
	collection := r.DB.Database("digital-voter").Collection("users")
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"password":           user.Password,
			"updated_at":         user.UpdatedAt,
			"reset_token":        nil,
			"reset_token_expiry": nil,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepository) UpdateProfiles(ctx context.Context, user *entity.User, oldEmail string) error {
	collection := r.DB.Database("digital-voter").Collection("users")

	filter := bson.M{"email": oldEmail}
	now := util.NowInWIB()
	user.UpdatedAt = &now
	update := bson.M{
		"$set": bson.M{
			"password":  user.Password,
			"email":     user.Email,
			"name":      user.Name,
			"updatedAt": user.UpdatedAt,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepository) UpdateScore(ctx context.Context, user *entity.User) error {
    collection := r.DB.Database("digital-voter").Collection("users")
    
    filter := bson.M{"email": user.Email}
    update := bson.M{
        "$set": bson.M{
            "score": user.Score,
            "updated_at": util.NowInWIB(),
        },
    }
    
    _, err := collection.UpdateOne(ctx, filter, update)
    return err
}
