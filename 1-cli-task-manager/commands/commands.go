package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

func InitStore() error {
	var err error
	db, err = bolt.Open("./tasks.db", 0600, nil)
	if err != nil {
		return err
	}
	// defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func CloseStore() error {
	return db.Close()
}

type Task struct {
	Id          int
	Description string
	Done        bool
}

func NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Action: func(cCtx *cli.Context) error {
			var err error
			var data []byte
			description := cCtx.Args().First()

			return db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("tasks"))
				fmt.Println("add task: ", description)
				id, _ := b.NextSequence()
				task := Task{
					Id:          int(id),
					Description: description,
					Done:        false,
				}
				data, err = json.Marshal(task)
				if err != nil {
					return err
				}
				return b.Put([]byte(strconv.Itoa(task.Id)), data)
			})
		},
	}
}

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "print task list",
		Action: func(cCtx *cli.Context) error {
			return db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("tasks"))
				c := b.Cursor()
				fmt.Println("You have the following tasks:")
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var task Task
					err := json.Unmarshal(v, &task)
					if err != nil {
						return err
					}
					if task.Done {
						fmt.Printf("%v. [x] %v\n", task.Id, task.Description)
					}
					if !task.Done {
						fmt.Printf("%v. [ ] %v\n", task.Id, task.Description)
					}
				}
				return nil
			})
		},
	}
}

func NewDoCommand() *cli.Command {
	return &cli.Command{
		Name:    "do",
		Aliases: []string{"d"},
		Usage:   "complete a task",
		Action: func(cCtx *cli.Context) error {
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("tasks"))
				id, err := strconv.Atoi(cCtx.Args().First())
				if err != nil {
					return err
				}
				v := b.Get([]byte(strconv.Itoa(id)))
				if v == nil {
					return errors.New("task not found")
				}
				task := Task{}
				err = json.Unmarshal(v, &task)
				if err != nil {
					return err
				}
				task.Done = true
				data, err := json.Marshal(task)
				if err != nil {
					return err
				}
				return b.Put([]byte(strconv.Itoa(task.Id)), data)
			})

			return nil
		},
	}
}
