package main

import (
	"github.com/go-martini/martini"
	"time"
	"github.com/martini-contrib/render"
	"code.google.com/p/leveldb-go/leveldb"
	"code.google.com/p/leveldb-go/leveldb/db"
	"strings"
	"log"
	"encoding/json"
	"github.com/martini-contrib/binding"
)

type Message struct {
	Id int             `form:"id"`
	Name string        `form:"name"`
	Mail string        `form:"mail"`
	Message string     `form:"message"`
	PostData time.Time
}

func main() {
	m := martini.Classic()

	m.Use(martini.Static("public"))

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
		Extensions: []string{".tmpl.html"},
	}))

	var opts db.Options
	level, err := leveldb.Open("database/messageDB", &opts)
	if err != nil {
		panic(err)
	}
	defer level.Close()

	m.Map(level)

	m.Get("/", Index)
	m.Get("/add", Add)
	m.Post("/create", binding.Bind(Message{}), Create)

	m.Run()
}

type ResponseData struct {
	Msg string
	Datas []Message
}

func Index(l *log.Logger, r render.Render, level *leveldb.DB) {
	// Message All
	var response ResponseData

	var readOpts db.ReadOptions
	keys, err := level.Get([]byte("keys"), &readOpts)
	if err == nil {
		keyList := strings.Split(string(keys), ",")

		response.Msg = "データベースのデータ"
		datas := make([]Message, len(keyList))
		for idx,key := range keyList {
			data, _ := level.Get([]byte(key), &readOpts)

			var message Message
			err := json.Unmarshal(data, &message)
			if err != nil {
				panic(err)
			}
			l.Println(key, message)
			datas[idx] = message
		}
		response.Datas = datas
	} else {
		// not database recodes
		response.Msg = "データベースのサンプル"
		response.Datas = []Message{
			{1, "ほげ", "test@test.com", "ふがふが", time.Now()},
			{2, "ねぷねぷ", "nepu@test.com", "ねぷてぬ", time.Now()},
		}
	}

	r.HTML(200, "index", response)
}
func Add(r render.Render) {
	var response ResponseData
	response.Msg = "入力してください"
	r.HTML(200, "add", response)
}

func Create(ms Message, level *leveldb.DB) (int, string) {

	// input Message
	inputData := ms
	inputData.PostData = time.Now()

	if inputData.Name == "" || inputData.Message == "" {
		return 400, "error"
	}

	var readOpts db.ReadOptions

	// inputData toJson
	inputJsonData, err := json.Marshal(inputData)
	if err != nil {
		panic(err)
	}

	// inputData save
	var writeOpts db.WriteOptions
	err = level.Set([]byte(inputData.Name), inputJsonData, &writeOpts)
	if err != nil {
		panic(err)
	}

	// find keys
	var lastedKeys string
	keys, err := level.Get([]byte("keys"), &readOpts)
	if err == nil {
		lastedKeys = string(keys)
		lastedKeys += ","
	} else {
		lastedKeys = ""
	}
	// update keys
	lastedKeys += inputData.Name
	err = level.Set([]byte("keys"), []byte(lastedKeys), &writeOpts)
	if err != nil {
		panic(err)
	}

	return 200, "ok"
}
