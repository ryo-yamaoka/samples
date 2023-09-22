package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const (
	usersTableName = "Users"
)

type User struct {
	UserID    string    `spanner:"UserID"`
	Name      string    `spanner:"Name"`
	CreatedAt time.Time `spanner:"CreatedAt"`
	UpdatedAt time.Time `spanner:"UpdatedAt"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", os.Getenv("SPANNER_PROJECT_ID"), os.Getenv("SPANNER_INSTANCE_ID"), os.Getenv("SPANNER_DATABASE_ID"))
	cli, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		return err
	}
	defer cli.Close()
	defer cleanup(cli)

	var mm []*spanner.Mutation

	// INSERT sample data
	pk := uuid.NewString()
	u := User{
		UserID:    pk,
		Name:      "aaa",
		CreatedAt: spanner.CommitTimestamp,
		UpdatedAt: spanner.CommitTimestamp,
	}
	m, err := spanner.InsertStruct(usersTableName, u)
	if err != nil {
		return err
	}
	mm = append(mm, m)
	createdAt, err := cli.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
		if err := rwt.BufferWrite(mm); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	// UPDATE Name after 3 seconds sleep
	_, err = cli.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
		row, err := rwt.ReadRow(ctx, usersTableName, spanner.Key{pk}, []string{"UserID", "Name", "CreatedAt", "UpdatedAt"})
		if err != nil {
			return err
		}
		var u User
		if err := row.ToStruct(&u); err != nil {
			return err
		}
		u.Name = "bbb"
		u.UpdatedAt = spanner.CommitTimestamp
		m, err := spanner.UpdateStruct(usersTableName, &u)
		if err != nil {
			return err
		}
		if err := rwt.BufferWrite([]*spanner.Mutation{m}); err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}

	// Read latest: expected name = 'bbb' and CreatedAt and UpdatedAt timestamps are difference
	fmt.Printf("StrongRead:\n")
	if err := read(ctx, cli, spanner.StrongRead()); err != nil {
		return err
	}

	// Read before UPDATE: expected name = 'aaa' and CreatedAt and UpdatedAt timestamps are the same
	fmt.Printf("Timestamp CreatedAt:\n")
	if err := read(ctx, cli, spanner.ReadTimestamp(createdAt)); err != nil {
		return err
	}

	// Read before CREATE within the maximum staleness timestamp: expected no record
	fmt.Printf("Timestamp too Old(-1h):\n")
	if err := read(ctx, cli, spanner.ReadTimestamp(createdAt.Add(-1*time.Hour))); err != nil {
		return err
	}

	// Read before CREATE with over the maximum staleness timestamp: expected FailedPrecondition
	fmt.Printf("Timestamp too Old(-1.5h):\n")
	if err := read(ctx, cli, spanner.ReadTimestamp(createdAt.Add(-1*time.Hour).Add(-30*time.Minute))); err != nil {
		return err
	}

	return nil
}

func read(ctx context.Context, cli *spanner.Client, tb spanner.TimestampBound) error {
	stmt := spanner.NewStatement("SELECT * FROM Users")
	iter := cli.ReadOnlyTransaction().WithTimestampBound(tb).Query(ctx, stmt)
	for {
		row, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return err
		}
		var u User
		if err := row.ToStruct(&u); err != nil {
			return err
		}
		fmt.Printf("%+v\n", u)
	}
	fmt.Println("")
	return nil
}

func cleanup(cli *spanner.Client) error {
	ctx := context.Background()
	if _, err := cli.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
		if err := rwt.BufferWrite([]*spanner.Mutation{spanner.Delete(usersTableName, spanner.AllKeys())}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
