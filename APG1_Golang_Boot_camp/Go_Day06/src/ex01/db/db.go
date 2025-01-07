package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Article struct {
	ID 			int
	Title 		string
	Content		string
}

type ArticleManagment interface {
	InsertArticleIntoDatabase(Article) error
	GetArticleByID(conn *pgx.Conn, id int) (*Article, error)
	GetArticlesList(conn *pgx.Conn) (*[]Article, error)
}

func main() {
	conn, err := NewPostgres("krillkovalev", "108814", "localhost", "5432", "my_blog")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer conn.Close(context.Background())
}

func NewPostgres(username, password, host, port, dbName string) (*pgx.Conn, error) {
    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil{
        return nil, err
    }
    return conn, nil
}

// InsertArticleIntoDatabase inserts a new article record into the database.
func InsertArticleIntoDatabase(conn *pgx.Conn, articleInfo Article) error {


	query := `
		insert into articles(id, title, content) values (@id, @title, @content) 
	`

	args := pgx.NamedArgs{
		"id":			articleInfo.ID,
		"title":		articleInfo.Title,
		"content":		articleInfo.Content,
	}

	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Inserting Article Into Database")
		return err
	}
	return nil
}

func GetArticleByID(conn *pgx.Conn, id int) (*Article, error) {

	query := `
		select * from articles where id = @articleID 
	`

	args := pgx.NamedArgs{
		"articleID": id,
	}

	row := conn.QueryRow(context.Background(), query, args)

	var article Article

	err := row.Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		log.Printf("Error Fetching Article Details")
		return nil, err
	}

	return &article, nil
}

func GetArticlesList(conn *pgx.Conn, page int) (*[]Article, error) { 

	limit := 3
	offset := limit * (page - 1)

	query := `
		select * from articles order by id limit @limit offset @offset
	`
	rows, err := conn.Query(context.Background(), query, limit, offset)
	if err != nil {
		log.Printf("Error Querying The Table")
		return nil, err
	}

	defer rows.Close()

	var articles []Article
	
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content)
		if err != nil {
			log.Printf("Error Fetching Articles")
			return &articles, err
		}
		articles = append(articles, article)
	}
	return &articles, nil
}