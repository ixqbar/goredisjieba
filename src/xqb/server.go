package xqb

import (
	"errors"
	"fmt"
	redis "github.com/jonnywang/go-kits/redis"
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
	VERSION = "0.0.1"
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

func (this *SearchRedisHandle) Init(db int) error {
	this.Lock()
	defer this.Unlock()

	dictPath := path.Join(jiebaXmlConfig.DictPath, fmt.Sprintf("%d", jiebaXmlConfig.DB))
	dictListPath := []string{
		fmt.Sprintf("%s/jieba.dict.utf8", dictPath),
		fmt.Sprintf("%s/hmm_model.utf8", dictPath),
		fmt.Sprintf("%s/user.dict.utf8", dictPath),
		fmt.Sprintf("%s/idf.utf8", dictPath),
		fmt.Sprintf("%s/stop_words.utf8", dictPath),
	}

	for _, p := range dictListPath {
		r, e := os.Stat(p)
		if e != nil || r.Size() == 0 {
			return errors.New(fmt.Sprintf("not found dict file `%s`", p))
		}
	}

	if this.jieba == nil {
		this.jieba = make(map[int]*gojieba.Jieba, 0)
	}

	this.jieba[db] = gojieba.NewJieba(dictListPath[0], dictListPath[1], dictListPath[2], dictListPath[3], dictListPath[4])

	return nil
}

func (this *SearchRedisHandle) Shutdown() error {
	redis.Logger.Print("searcher server will shutdown!!!")

	if this.jieba != nil {
		for _, j := range this.jieba {
			j.Free()
		}
	}

	return nil
}

func (this *SearchRedisHandle) Version() (string, error) {
	return VERSION, nil
}

func (this *SearchRedisHandle) Select(client *redis.Client, db int) (string, error) {
	err := this.Init(db)
	if err != nil {
		return "", err
	}

	client.DB = db

	return OK, nil
}

func (this *SearchRedisHandle) CutAll(client *redis.Client, words string) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return this.jieba[client.DB].CutAll(words), nil
}

func (this *SearchRedisHandle) Cut(client *redis.Client, words string, useHmm int) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return this.jieba[client.DB].Cut(words, toBool(useHmm)), nil
}

func (this *SearchRedisHandle) CutForSearch(client *redis.Client, words string, useHmm int) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return this.jieba[client.DB].CutForSearch(words, toBool(useHmm)), nil
}

func (this *SearchRedisHandle) Tag(client *redis.Client, words string) ([]string, error) {
	if len(words) == 0 {
		return nil, ERRPARAMS
	}

	return this.jieba[client.DB].Tag(words), nil
}

func (this *SearchRedisHandle) AddWord(client *redis.Client, word string) (string, error) {
	if len(word) == 0 {
		return "", ERRPARAMS
	}

	this.jieba[client.DB].AddWord(word)

	return OK, nil
}

func Run() {
	searcher := &SearchRedisHandle{
		jieba: nil,
	}

	searcher.SetShield("Init")
	searcher.SetShield("Shutdown")
	searcher.SetShield("Lock")
	searcher.SetShield("Unlock")
	searcher.SetShield("SetShield")
	searcher.SetShield("SetConfig")
	searcher.SetShield("CheckShield")

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
