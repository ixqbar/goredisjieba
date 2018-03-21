package xqb

import (
	"errors"
	"fmt"
	"github.com/jonnywang/go-kits/redis"
	"github.com/yanyiwu/gojieba"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

var (
	ERRPARAMS = errors.New("error params")
)

const (
	VERSION = "0.0.3"
	OK      = "OK"
)

type SearchRedisHandle struct {
	redis.RedisHandler
	sync.Mutex
	jieba map[int]*gojieba.Jieba
}

func toBool(p int) bool {
	if p == 0 {
		return false
	}

	return true
}

func (obj *SearchRedisHandle) Init(db int) error {
	obj.Lock()
	defer obj.Unlock()

	if _, ok := obj.jieba[db]; ok {
		return nil
	}

	dictPath := path.Join(jiebaXmlConfig.DictPath, fmt.Sprintf("%d", db))
	dictPaths := []string{
		fmt.Sprintf("%s/jieba.dict.utf8", dictPath),
		fmt.Sprintf("%s/hmm_model.utf8", dictPath),
		fmt.Sprintf("%s/user.dict.utf8", dictPath),
		fmt.Sprintf("%s/idf.utf8", dictPath),
		fmt.Sprintf("%s/stop_words.utf8", dictPath),
	}

	for _, p := range dictPaths {
		redis.Logger.Print(p)
		r, e := os.Stat(p)
		if e != nil || r.Size() == 0 {
			return errors.New(fmt.Sprintf("not found dict file `%s`", p))
		}
	}

	obj.jieba[db] = gojieba.NewJieba(dictPaths[0], dictPaths[1], dictPaths[2], dictPaths[3], dictPaths[4])

	return nil
}

func (obj *SearchRedisHandle) Shutdown() error {
	redis.Logger.Print("searcher server will shutdown!!!")

	if obj.jieba != nil {
		for _, j := range obj.jieba {
			j.Free()
		}
	}

	return nil
}

func (obj *SearchRedisHandle) Version() (string, error) {
	return VERSION, nil
}

func (obj *SearchRedisHandle) Ping(content string) (string, error)  {
	if len(content) > 0 {
		return content, nil
	}

	return "PONG", nil
}

func (obj *SearchRedisHandle) Select(client *redis.Client, db int) (string, error) {
	err := obj.Init(db)
	if err != nil {
		return "", err
	}

	client.DB = db

	return OK, nil
}

func (obj *SearchRedisHandle) CutAll(client *redis.Client, words string) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return obj.jieba[client.DB].CutAll(words), nil
}

func (obj *SearchRedisHandle) Cut(client *redis.Client, words string, useHmm int) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return obj.jieba[client.DB].Cut(words, toBool(useHmm)), nil
}

func (obj *SearchRedisHandle) CutForSearch(client *redis.Client, words string, useHmm int) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return obj.jieba[client.DB].CutForSearch(words, toBool(useHmm)), nil
}

func (obj *SearchRedisHandle) Tag(client *redis.Client, words string) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return obj.jieba[client.DB].Tag(words), nil
}

func (obj *SearchRedisHandle) Extract(client *redis.Client, words string, limit int) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return obj.jieba[client.DB].Extract(words, limit), nil
}

func (obj *SearchRedisHandle) AddWord(client *redis.Client, word string) (string, error) {
	if len(word) == 0 {
		return "", ERRPARAMS
	}

	obj.jieba[client.DB].AddWord(word)

	return OK, nil
}

func Run() {
	searcher := &SearchRedisHandle{
		jieba: make(map[int]*gojieba.Jieba, 0),
	}

	searcher.Initiation(nil)

	err := searcher.Init(jiebaXmlConfig.DB)
	if err != nil {
		redis.Logger.Print(err)
		return
	}

	server, err := redis.NewServer(jiebaXmlConfig.Address, searcher)
	if err != nil {
		redis.Logger.Print(err)
		return
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		server.Stop(10)
	}()

	redis.Logger.Printf("server run at %s", jiebaXmlConfig.Address)

	err = server.Start()
	if err != nil {
		redis.Logger.Print(err)
	}

	searcher.Shutdown()
}
