// Copyright 2019 ihciah <ihciah@gmail.com>
//
// Licensed under the GNU Affero General Public License, Version 3.0
// (the "License"); you may not use this file except in compliance with the
// License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Requester struct {
	Client          http.Client
	header          map[string]string
	reqId           int64
	lastRequestTime int64
	lock            sync.Mutex
	withProxy       bool
	lastProxyTime   int64
	RequestLimiter  *RequestLimiter
	logger          *log.Logger
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func (kv *KeyValue) Marshal() string {
	var valuestr string
	switch t := kv.Value.(type) {
	case int64:
		valuestr = strconv.FormatInt(t, 10)
	case string:
		valuestr = t
	case int:
		valuestr = strconv.Itoa(t)
	}
	s := url.QueryEscape(valuestr)
	s = strings.Replace(s, "%2F", "%2f", -1)
	s = strings.Replace(s, "%2B", "%2b", -1)
	s = strings.Replace(s, "%3D", "%3d", -1)
	s = strings.Replace(s, "%2D", "%2d", -1)
	return fmt.Sprintf("%s=%s", kv.Key, url.QueryEscape(valuestr))
}

type OrderedMap struct {
	kv []KeyValue
}

func MakeOrderedMap() OrderedMap {
	return OrderedMap{make([]KeyValue, 0, initOrderedMapSize)}
}

func (om *OrderedMap) Add(key string, value interface{}) {
	om.kv = append(om.kv, KeyValue{key, value})
}

func (om *OrderedMap) LoadMap(m map[string]interface{}) {
	for k, v := range m {
		om.Add(k, v)
	}
}

func (om *OrderedMap) Marshal() string {
	tobeconcat := make([]string, 0, len(om.kv))
	for _, kv := range om.kv {
		tobeconcat = append(tobeconcat, kv.Marshal())
	}
	return strings.Join(tobeconcat, "&")
}

func GetProxy() (*url.URL, error) {
	resp, err := http.Get(proxyAPI)
	if err != nil {
		log.Println("Proxy get error!", err)
		return nil, ErrProxyApiUnavailable
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	proxyUrl, err := url.Parse("http://" + strings.Trim(string(body), " \n\r\t,"))
	log.Println("Use proxy:", proxyUrl, "Error:", err)
	if err == nil {
		return proxyUrl, nil
	}
	return proxyUrl, ErrProxyApiUnavailable
}

func MakeRequester(channelId string, RequestLimiter *RequestLimiter, withProxy, debugger bool) Requester {
	headerAndroid := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Dalvik/1.6.0 (Linux; U; Android 4.4.2; MI 6  Build/NMF26X)",
	}
	headerIos := map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"X-Unity-Version": "5.2.5f1",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-cn",
		"Accept":          "*/*",
		"User-Agent":      "girlsfrontline/563 CFNetwork/902.2 Darwin/17.7.0",
	}
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar:     cookieJar,
		Timeout: time.Duration(networkTimeout * time.Second),
	}

	if debugger {
		proxyUrl, _ := url.Parse("http://127.0.0.1:9999")
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	header := headerIos
	if channelId != "ios" {
		header = headerAndroid
	}

	return Requester{Client: client, header: header, reqId: time.Now().Unix() * 100000, RequestLimiter: RequestLimiter,
		withProxy: withProxy, logger: log.New(os.Stdout, "[Requester]", log.LstdFlags)}
}

func (r *Requester) RefreshProxy() {
	if r.withProxy {
		now := time.Now().Unix()
		sleep := proxyInterval - now + r.lastProxyTime
		if sleep > 0 {
			time.Sleep(time.Duration(sleep) * time.Second)
		}
		proxyUrl, err := GetProxy()
		r.lastProxyTime = time.Now().Unix()
		if err == nil {
			r.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		}
	}
}

func (r *Requester) ReadRes(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.header {
		req.Header.Set(k, v)
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		r.logger.Println("ERR when send request", err)
		return nil, ErrNetworkUnavailable
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.logger.Println("ERR when read response", err)
		return respBody, ErrNetworkUnavailable
	}
	if resp.StatusCode != 200 {
		r.logger.Println("Status code:", resp.StatusCode)
		return respBody, ErrHttpCodeError
	}
	return respBody, err
}

func (r *Requester) Get(url string, limit bool) ([]byte, error) {
	if limit && r.RequestLimiter != nil {
		r.RequestLimiter.Wait()
	}
	for errCount := 0; errCount < maxErrorCount; errCount++ {
		resp, err := r.ReadRes("GET", url, nil)
		if err == ErrNetworkUnavailable {
			r.logger.Print("Network error: ", err)
			time.Sleep(networkErrorWait * time.Second)
		} else {
			return resp, err
		}
	}
	return []byte{}, ErrTooManyFailure
}

func (r *Requester) Post(url string, body []byte, limit bool) ([]byte, error) {
	if limit && r.RequestLimiter != nil {
		// 全局请求速率控制
		r.RequestLimiter.Wait()
	}
	for errCount := 0; errCount < maxErrorCount; errCount++ {
		resp, err := r.ReadRes("POST", url, bytes.NewBuffer(body))
		if err == ErrNetworkUnavailable {
			r.logger.Print("Network error: ", err)
			respStr := string(resp)
			r.logger.Print("Response: ", respStr)
			if strings.HasPrefix(respStr, "{") && strings.HasSuffix(respStr, "}") {
				return resp, nil
			}
			time.Sleep(networkErrorWait * time.Second)
		} else if err != nil {
			r.logger.Print("Response: ", string(resp))
			return resp, err
		} else {
			return resp, err
		}
	}
	return []byte{}, ErrTooManyFailure
}

func (r *Requester) GetReqID() int64 {
	reqId := atomic.AddInt64(&r.reqId, 1)
	r.lock.Lock()
	defer r.lock.Unlock()

	// 单用户的请求速率控制
	timeInterval := r.lastRequestTime + requestIntervalNano - time.Now().UnixNano()
	if timeInterval > 0 {
		time.Sleep(time.Duration(timeInterval) * time.Nanosecond)
	}
	r.lastRequestTime = time.Now().UnixNano()

	return reqId
}

const (
	normalMode = iota
	maintainMode
)

type RequestLimiter struct {
	Interval         int64
	MaintainInterval int64
	last             int64
	lock             sync.Mutex
	statusLock       sync.RWMutex
	status           int
	waitBreaker      chan bool
	logger           *log.Logger
}

func (rl *RequestLimiter) Wait() {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	var interval int64
	rl.statusLock.RLock()
	if rl.status == normalMode {
		interval = rl.Interval
	} else {
		interval = rl.MaintainInterval
	}
	rl.statusLock.RUnlock()

	toSleep := rl.last + interval - time.Now().Unix()
	if toSleep > 0 {
		rl.logger.Printf("RequestLimiter will sleep %d seconds.", toSleep)
		select {
		case <-time.After(time.Duration(toSleep) * time.Second):
			break
		case <-rl.waitBreaker:
			break
		}
	}
	rl.last = time.Now().Unix()
}

// 游戏服务器进入维护模式，这时将该limiter间隔临时改为MaintainInterval。
func (rl *RequestLimiter) EnterMaintainMode() {
	rl.statusLock.RLock()
	statusNow := rl.status
	rl.statusLock.RUnlock()

	if statusNow == maintainMode {
		return
	}

	rl.statusLock.Lock()
	defer rl.statusLock.Unlock()
	rl.status = maintainMode
	rl.logger.Print("Enter maintain mode.")
}

// 游戏服务器退出维护模式，这时将恢复limiter间隔，同时唤醒正在sleep的Wait。
func (rl *RequestLimiter) ExitMaintainMode() {
	rl.statusLock.RLock()
	statusNow := rl.status
	rl.statusLock.RUnlock()

	if statusNow == normalMode {
		return
	}

	rl.statusLock.Lock()
	defer rl.statusLock.Unlock()
	rl.status = normalMode
	rl.logger.Print("Exit maintain mode.")
	select {
	case rl.waitBreaker <- true:
		rl.logger.Print("Break sleep request.")
		break
	default:
		break
	}
}

func MakeRequestLimiter(interval, maintainInterval int64) RequestLimiter {
	return RequestLimiter{Interval: interval, MaintainInterval: maintainInterval, waitBreaker: make(chan bool), logger: log.New(os.Stdout, "[RequestLimiter]", log.LstdFlags)}
}
