// Copyright 2022 The GPostgres Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gpostgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/alpstable/gidari/proto"
)

const defaultConnectionString = "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"

func testClient(t *testing.T, ctx context.Context, connString string) *sql.DB {
	t.Helper()

	database, err := sql.Open("postgres", connString)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	return database
}

func defaultTestClient(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	return testClient(t, ctx, defaultConnectionString)
}

func TestPostgres(t *testing.T) {
	t.Parallel()

	defaultTestTable := &proto.Table{Name: "tests1"}
	listTablesTable := &proto.Table{Name: "lttests1"}
	listPrimaryKeysTable := &proto.Table{Name: "pktests1"}

	defaultData := map[string]interface{}{
		"test_string": "test",
		"id":          "1",
	}

	ctx := context.Background()

	pg, err := New(ctx, defaultTestClient(t, ctx))
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	proto.RunTest(context.Background(), t, pg, func(runner *proto.TestRunner) {
		runner.AddCloseDBCases(
			[]proto.TestCase{
				{
					Name: "close postgres",
					OpenFn: func() proto.Storage {
						stg, _ := New(ctx, defaultTestClient(t, ctx))
						return stg
					},
				},
			}...,
		)

		runner.AddStorageTypeCases(
			[]proto.TestCase{
				{
					Name:        "storage type",
					StorageType: proto.PostgresType,
				},
			}...,
		)

		runner.AddIsNoSQLCases(
			[]proto.TestCase{
				{
					Name:            "isNoSQL postgres",
					ExpectedIsNoSQL: false,
				},
			}...,
		)

		runner.AddListPrimaryKeysCases(
			[]proto.TestCase{
				{
					Name:  "single",
					Table: listPrimaryKeysTable,
					ExpectedPrimaryKeys: map[string][]string{
						listPrimaryKeysTable.Name: {"test_string"},
					},
				},
			}...,
		)

		runner.AddListTablesCases(
			[]proto.TestCase{
				{
					Name:  "single",
					Table: listTablesTable,
				},
			}...,
		)

		runner.AddUpsertTxnCases(
			[]proto.TestCase{
				{
					Name:               "commit",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 8192,
					Data:               defaultData,
				},
				{
					Name:               "rollback",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 0,
					Rollback:           true,
					Data:               defaultData,
				},
				{
					Name:               "rollback on error",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 0,
					ForceError:         true,
					Data:               defaultData,
				},
			}...,
		)

		runner.AddUpsertBinaryCases(
			[]proto.TestCase{
				{
					Name:               "no pk map",
					BinaryColumn:       "data",
					Table:              &proto.Table{Name: "property_bag_tests1"},
					ExpectedUpsertSize: 8192,
					Data: map[string]interface{}{
						"data": []byte("{ x: 1 }"),
						"id":   "1",
					},
				},
				{
					Name:         "pk map",
					BinaryColumn: "data",
					Table:        &proto.Table{Name: "property_bag_tests2"},
					PrimaryKeyMap: map[string]string{
						"pk1": "primary_key1",
						"pk2": "primary_key2",
					},
					ExpectedUpsertSize: 8192,
					Data: map[string]interface{}{
						"data": []byte("{ x: 1 }"),
						"pk1":  "1",
						"pk2":  "2",
					},
				},
			}...,
		)

		runner.AddPingCases(
			[]proto.TestCase{
				{
					Name: "check postgres connection",
				},
			}...,
		)

		runner.AddTruncateCases(
			[]proto.TestCase{
				{
					Name: "Trucate case",
					OpenFn: func() proto.Storage {
						stg, _ := New(ctx, defaultTestClient(t, ctx))
						return stg
					},
				},
			}...,
		)
	})
}
