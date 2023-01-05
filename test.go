package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	//"github.com/docker/docker/client"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/test", testHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	//TODO remove testing
	//fmt.Println("port:", port)

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//	defer cancel()

	//Create database, tables
	//adminClient, err12 := database.NewDatabaseAdminClient(ctx)
	//if err12 != nil {
	//	log.Printf("err12", err12)
	//}
	//
	//op, err10 := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
	//	Parent:          "projects/rp-test-s-01/instances/test-instance",
	//	CreateStatement: "CREATE DATABASE `example-db`",
	//	ExtraStatements: []string{
	//		`CREATE TABLE Singers (
	//			SingerId   INT64 NOT NULL,
	//			FirstName  STRING(1024),
	//			LastName   STRING(1024),
	//			SingerInfo BYTES(MAX)
	//		) PRIMARY KEY (SingerId)`,
	//		`CREATE TABLE Albums (
	//			SingerId     INT64 NOT NULL,
	//			AlbumId      INT64 NOT NULL,
	//			AlbumTitle   STRING(MAX)
	//		) PRIMARY KEY (SingerId, AlbumId),
	//		INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
	//	},
	//})
	//
	//if err10 != nil {
	//	log.Printf("err10", err10)
	//}
	//
	//if _, err11 := op.Wait(ctx); err11 != nil {
	//	log.Printf("err11", err11)
	//}
	//fmt.Fprintf(w, "Created database ")

	// Client
	dataClient, err1 := spanner.NewClient(ctx, "projects/rp-test-s-01/instances/test-instance/databases/example-db")
	if err1 != nil {
		log.Fatal(err1)
	}
	//	defer dataClient.Close()

	//DML insert multiple rows
	_, err2 := dataClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName) VALUES
				(12, 'Melissa', 'Garcia'),
				(13, 'Russell', 'Morales'),
				(14, 'Jacqueline', 'Long'),
				(15, 'Dylan', 'Shaw')`,
		}
		rowCount, err3 := txn.Update(ctx, stmt)
		if err3 != nil {
			return err3
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err3
	})
	log.Printf("err2", err2)

	// mutation
	//singerColumns := []string{"SingerId", "FirstName", "LastName"}
	//m := []*spanner.Mutation{
	//	spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{1, "Marc", "Richards"}),
	//	spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{2, "Catalina", "Smith"}),
	//	spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{3, "Alice", "Trentor"}),
	//	spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{4, "Lea", "Martin"}),
	//	spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{5, "David", "Lomond"}),
	//}
	//_, err21 := dataClient.Apply(ctx, m)
	//if err21 != nil {
	//	log.Printf("err21", err21)
	//}
	//log.Printf("err21", err21)

	// mutation test - slice of structs - submitted a support case on this, waiting for response
	//type name struct {
	//	ID        int64
	//	FirstName string
	//	LastName  string
	//}
	//var singers []name
	//
	//var singerInfo = name{1, "Timothy", "Campbell"}
	//singers = append(singers, singerInfo)
	//singerInfo = name{2, "Timothy2", "Campbell2"}
	//singers = append(singers, singerInfo)
	//
	//singerColumns := []string{"SingerId", "FirstName", "LastName"}
	//m := []*spanner.Mutation{
	//	spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{singerInfo.ID, singerInfo.FirstName, singerInfo.LastName}),
	//}
	//_, err21 := dataClient.Apply(ctx, m)
	//if err21 != nil {
	//	log.Printf("err21", err21)
	//}
	//log.Printf("err21", err21)

	//Query with parameter
	//stmt := spanner.Statement{
	//	SQL: `SELECT SingerId, FirstName, LastName FROM Singers
	//		WHERE SingerId > @singerId`,
	//	Params: map[string]interface{}{
	//		"singerId": 1,
	//	},
	//}
	//iter := dataClient.Single().Query(ctx, stmt)
	//defer iter.Stop()
	//for {
	//	row, err31 := iter.Next()
	//	if err31 == iterator.Done {
	//		return
	//	}
	//	if err31 != nil {
	//		log.Printf("err31", err31)
	//	}
	//	var singerID int64
	//	var firstName, lastName string
	//	if err32 := row.Columns(&singerID, &firstName, &lastName); err32 != nil {
	//		log.Printf("err32", err32)
	//		return
	//	}
	//	log.Printf("singerID", singerID, firstName, lastName)
	//	fmt.Fprintf(w, "%d %s %s\n", singerID, firstName, lastName)
	//}

	// DML using struct - doesn't look like can iterate over a slice of structs
	//_, err41 := dataClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
	//	type name struct {
	//		FirstName string
	//		LastName  string
	//	}
	//	var singerInfo = name{"Timothy", "Campbell"}
	//
	//	stmt := spanner.Statement{
	//		SQL: `Update Singers Set LastName = 'Grant'
	//			WHERE STRUCT<FirstName String, LastName String>(Firstname, LastName) = @name`,
	//		Params: map[string]interface{}{"name": singerInfo},
	//	}
	//	rowCount, err := txn.Update(ctx, stmt)
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
	//	return nil
	//})
	//if err41 != nil {
	//	log.Printf("err41", err41)
	//}

}
