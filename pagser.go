// Copyright 2020 Foolin

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pagser

import "errors"

// Pagser the page parser
type Pagser struct {
	Config   Config
	ctxTags  map[string]*parseTag // tag value => parseTag
	ctxFuncs map[string]CallFunc  // name => func
}

// New create client
func New() *Pagser {
	p, _ := NewWithConfig(DefaultConfig())
	return p
}

// NewWithConfig create client with Config and error
func NewWithConfig(cfg Config) (*Pagser, error) {
	if cfg.TagerName == "" {
		return nil, errors.New("TagerName must not empty")
	}
	if cfg.FuncSymbol == "" {
		return nil, errors.New("FuncSymbol must not empty")
	}
	return &Pagser{
		Config:   cfg,
		ctxTags:  make(map[string]*parseTag, 0),
		ctxFuncs: builtinFuncMap,
	}, nil
}

// DefaultConfig the default Config
func DefaultConfig() Config {
	return defaultCfg
}
