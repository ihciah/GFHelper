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

import "errors"

var ErrServerError = errors.New("server return error")
var ErrUnknownError = errors.New("unknown error")
var ErrInputError = errors.New("input data format error")
var ErrWrongPass = errors.New("wrong username or password")
var ErrUserNotExist = errors.New("user does not exist")
var ErrWrongToken = errors.New("wrong openid or access token")
var ErrRegError = errors.New("error when register")
var ErrRegUsed = errors.New("email has benn used")
var ErrHttpCodeError = errors.New("wrong http code")
var ErrGame203 = errors.New("gun maximum exceeded(err203)")
var ErrGame3 = errors.New("doing action before mission started(err3)")
var ErrGame2 = errors.New("doing action before mission started(err2)")
var ErrGame = errors.New("error duing game")
var ErrServerMaintain = errors.New("server maintain")
var ErrProxyApiUnavailable = errors.New("proxy api error")
var ErrNetworkUnavailable = errors.New("network unavailable")
var ErrTooManyFailure = errors.New("too many failure")

const defaultKey = "yundoudou"
const maxRetire = 36
const IOSAppId = "0001000100021001"
const AndroidAppId = "0002000100021001"

const initOrderedMapSize = 10
const requestIntervalNano = 500000000
const networkErrorWait = 1
const maxErrorCount = 3
const networkTimeout = 3

const proxyAPI = "http://piping.mogumiao.com/proxy/api/get_ip_bs?appKey=&count=1&expiryDate=0&format=2&newLine=1"
const proxyInterval = 5
