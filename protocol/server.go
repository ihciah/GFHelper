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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Server struct {
	Addr      string
	Name      string
	Condition int
	IsCheck   int
	ClientVer string
	AbVer     string
	Version   Version
}

type Version struct {
	Now               string `json:"now"`
	TomorrowZero      string `json:"tomorrow_zero"`
	MonthZero         int    `json:"month_zero"`
	NextMonthZero     int    `json:"next_month_zero"`
	Timezone          string `json:"timezone"`
	DataVersion       string `json:"data_version"`
	ClientVersion     string `json:"client_version"`
	AbVersion         string `json:"ab_version"`
	Weekday           int    `json:"weekday"`
	AuthenticationURL string `json:"authentication_url"`
}

type serversXML struct {
	Servers []serverXML `xml:"server"`
	Config  configXML   `xml:"config"`
}

type serverXML struct {
	ServerName      string `xml:"name"`
	ServerAddr      string `xml:"addr"`
	ServerCondition int    `xml:"condition"`
	ServerIsCheck   int    `xml:"is_check"`
}

type configXML struct {
	ClientVer string `xml:"client_version"`
	AbVer     string `xml:"ab_version"`
}

// 获取服务器区域(取第一个)，并返回其状态和版本信息
func GetServer(channelid string, checkver string) (Server, error) {
	var server Server
	err := UpdateServer(channelid, checkver, &server)
	return server, err
}

// 获取服务器区域和对应Version信息，写入dst
func UpdateServer(channelid string, checkver string, dst *Server) error {
	var channel string
	var platformChannelId string
	var url string

	if channelid == "android" {
		channel = "cn_mica"
		platformChannelId = "GWPZ"
		url = "http://adr.transit.gf.ppgame.com/index.php"
	} else if channelid == "ios" {
		channel = "cn_appstore"
		platformChannelId = "ios"
		url = "http://ios.transit.gf.ppgame.com/index.php"
	} else if channelid == "tw" {
		channel = "as_tianxia"
		platformChannelId = "google"
		url = "http://sn-list.txwy.tw/index.php"
	}

	rnder := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := fmt.Sprintf("%06v", rnder.Int31n(1000000))

	strdata := fmt.Sprintf("c=game&a=newserverList&channel=%s&platformChannelId=%s&check_version=%s&rnd=%s",
		channel, platformChannelId, checkver, rnd)
	data := []byte(strdata)
	requester := MakeRequester(channelid, nil, false, false)
	ret, err := requester.Post(url, data, false)
	if err != nil {
		return err
	}
	var servers serversXML
	err = xml.Unmarshal(ret, &servers)
	if err != nil {
		return err
	}
	if len(servers.Servers) < 1 {
		log.Println("Can not get a server.")
		return ErrServerError
	}
	dst.Addr = servers.Servers[0].ServerAddr
	dst.Name = servers.Servers[0].ServerName
	dst.Condition = servers.Servers[0].ServerCondition
	dst.IsCheck = servers.Servers[0].ServerIsCheck
	dst.ClientVer = servers.Config.ClientVer
	dst.AbVer = servers.Config.AbVer

	//if dst.Condition != 0 {
	//	log.Println("Server maintain.")
	//	return ErrServerMaintain
	//}

	postdata := fmt.Sprintf("=&req_id=%d", requester.GetReqID())
	ret, err = requester.Post(dst.Addr+"Index/version", []byte(postdata), false)
	if err != nil {
		return err
	}
	err = json.Unmarshal(ret, &dst.Version)
	return err
}
