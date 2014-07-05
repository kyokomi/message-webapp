package main

import (
	"github.com/go-martini/martini"
	"time"
	"github.com/martini-contrib/render"
)


/*
@Entity
public class Message extends Model {

	@Id
	public Long id;
	public String name;
	public String mail;
	public String message;
	public Date postdate;

	public static Finder<Long, Message> find =
		new Finder<Long, Message>(Long.class, Message.class);

	@Override
	public String toString(){
		return ("[id:" + id + ", name:" + name + ", mail:" + mail +
			", message:" + message + ", date:" + postdate + "]");
	}
}
 */
type Message struct {
	id int
	name string
	mail string
	message string
	postData time.Time
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", Index)
	m.Get("/add", Add)
	m.Post("/create", Create)

	m.Run()
}

func Index(params martini.Params, r render.Render) {
	// Message All
//	return 200, "データベースのサンプル"
	r.HTML(200, "index", "test")
}
func Add(params martini.Params) (int, string) {
	return 200, "投稿フォーム"
}

func Create(params martini.Params) (int, string) {
	// input Message
	// database save
	return 200, "hoge\n"
}
