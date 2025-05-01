package mongodb

import (
	"app/pkg/database"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	*mongo.Client
	cfg *database.DatabaseConfig
}

func NewClient(cfg *database.DatabaseConfig) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-256",
		c.cfg.Username,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.Database,
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	c.Client = client
	return nil
}

func (c *Client) Disconnect() error {
	if c.Client != nil {
		return c.Client.Disconnect(context.Background())
	}
	return nil
}

func (c *Client) Ping() error {
	if c.Client != nil {
		return c.Client.Ping(context.Background(), nil)
	}
	return fmt.Errorf("mongodb connection not established")
}
