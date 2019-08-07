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

package gameact

import "errors"

var ErrUnexpectedData = errors.New("unexpected data")
var ErrNotLogin = errors.New("not login")
var ErrResourcesInsufficient = errors.New("resources insufficient")
var ErrNoFreeSpace = errors.New("no free space for new guns")

const WithSupply = false
const TimePadding = 20
const MaxRegFail = 1
const MaxLoginFail = 3
const MinMRE = 1000
const MinAMMO = 1000
const MinMP = 1000
const RetireGap = 5
const MinEatCount = 8
