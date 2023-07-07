package bigtable

import (
	"cloud.google.com/go/bigtable"
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
)

const (
	idColumnFamily     = "id"
	relColumnFamily    = "rel"
	csAuthColumnFamily = "cs_auth"
	typeColumn         = "t"
)

type Store struct {
	client *bigtable.Client
	table  *bigtable.Table
}

func NewStore(ctx context.Context, project, instance string) (*Store, error) {
	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		return nil, fmt.Errorf("bigtable admin client for %s/%s: %w", project, instance, err)
	}

	tables, err := adminClient.Tables(ctx)
	if err != nil {
		return nil, fmt.Errorf("bigtable tables for %s/%s: %w", project, instance, err)
	}

	const tableName = "manager"
	if !slices.Contains(tables, tableName) {
		log.Printf("creating table %s", tableName)
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			return nil, fmt.Errorf("bigtable create table %s for %s/%s: %w", tableName, project, instance, err)
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("bigtable table %s info for %s/%s: %w", tableName, project, instance, err)
	}

	columnFamilies := []string{idColumnFamily, relColumnFamily, csAuthColumnFamily}
	for _, columnFamily := range columnFamilies {
		if !slices.Contains(tblInfo.Families, columnFamily) {
			log.Printf("creating column family %s/%s", tableName, columnFamily)
			if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamily); err != nil {
				return nil, fmt.Errorf("bigtable %s create column family %s for %s/%s: %w", tableName, columnFamily, project, instance, err)
			}
		}
	}

	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		return nil, fmt.Errorf("bigtable client for %s/%s: %w", project, instance, err)
	}

	table := client.Open(tableName)

	return &Store{
		client: client,
		table:  table,
	}, nil
}
