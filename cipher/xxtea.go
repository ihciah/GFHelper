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

package cipher

const delta = 0x9E3779B9

func xxteaToBytes(v []uint32) []byte {
	length := uint32(len(v))
	n := length << 2
	bytes := make([]byte, n)
	for i := uint32(0); i < n; i++ {
		bytes[i] = byte(v[i>>2] >> ((i & 3) << 3))
	}
	return bytes
}

func xxteaToUint32s(bytes []byte) (v []uint32) {
	length := uint32(len(bytes))
	n := length >> 2
	if length&3 != 0 {
		n++
	}
	v = make([]uint32, n)
	for i := uint32(0); i < length; i++ {
		v[i>>2] |= uint32(bytes[i]) << ((i & 3) << 3)
	}
	return v
}

func xxteaMx(sum uint32, y uint32, z uint32, p uint32, e uint32, k []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (k[p&3^e] ^ z))
}

func xxteaFixk(k []uint32) []uint32 {
	if len(k) < 4 {
		key := make([]uint32, 4)
		copy(key, k)
		return key
	}
	return k
}

func xxteaEncrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = xxteaFixk(k)
	var y, z, sum, e, p, q uint32
	z = v[n]
	sum = 0
	for q = 6 + 52/length; q > 0; q-- {
		sum += delta
		e = sum >> 2 & 3
		for p = 0; p < n; p++ {
			y = v[p+1]
			v[p] += xxteaMx(sum, y, z, p, e, k)
			z = v[p]
		}
		y = v[0]
		v[n] += xxteaMx(sum, y, z, p, e, k)
		z = v[n]
	}
	return v
}

func xxteaDecrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = xxteaFixk(k)
	var y, z, sum, e, p, q uint32
	y = v[0]
	q = 6 + 52/length
	for sum = q * delta; sum != 0; sum -= delta {
		e = sum >> 2 & 3
		for p = n; p > 0; p-- {
			z = v[p-1]
			v[p] -= xxteaMx(sum, y, z, p, e, k)
			y = v[p]
		}
		z = v[n]
		v[0] -= xxteaMx(sum, y, z, p, e, k)
		y = v[0]
	}
	return v
}

func XXTEAEncrypt(data []byte) []byte {
	if data == nil || len(data) == 0 {
		return data
	}
	xxtea_key := []byte("1234567890\x00\x00\x00\x00\x00\x00")
	return xxteaToBytes(xxteaEncrypt(xxteaToUint32s(data), xxteaToUint32s(xxtea_key)))
}

func XXTEADecrypt(data []byte) []byte {
	if data == nil || len(data) == 0 {
		return data
	}
	xxtea_key := []byte("1234567890\x00\x00\x00\x00\x00\x00")
	return xxteaToBytes(xxteaDecrypt(xxteaToUint32s(data), xxteaToUint32s(xxtea_key)))
}

func XXTEAEncryptFromString(str string) []byte {
	return XXTEAEncrypt([]byte(str))
}

func XXTEADecryptToString(b []byte) (string, error) {
	result := XXTEADecrypt(b)
	return string(result), nil
}
