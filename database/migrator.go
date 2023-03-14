package database

import (
	"context"
	"github.com/frchandra/chatin/database/factory"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migrator struct {
	db  *mongo.Database
	log *logrus.Logger
}

func NewMigrator(db *mongo.Database, log *logrus.Logger) *Migrator {
	return &Migrator{db: db, log: log}
}

func (m *Migrator) RunMigration() {
	if err := m.db.Collection("users").Drop(context.Background()); err != nil {
		m.log.Warn("error when dropping collection. Error: ", err.Error())
	}
	if err := m.db.CreateCollection(context.Background(), "users"); err != nil {
		m.log.Error("error when creating collection. Error: " + err.Error())
	}
	if err := m.RunSeeder(); err != nil {
		m.log.Error("error when running seeder. Error: " + err.Error())
	}
	m.log.Info("migration runned successfully")
}

func (m *Migrator) GetFactory() []factory.Factory {
	return []factory.Factory{
		factory.NewUserFactory(m.db, m.log),
	}
}

func (m *Migrator) RunSeeder() error {
	for _, seeder := range m.GetFactory() {
		if err := seeder.RunFactory(); err != nil {
			return err
		}
	}
	return nil
}
