// Copyright 2020 Foolin

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pagser

import "errors"

type Pagser struct {
	config Config
	tagers map[string]*Tager   // tager value => Tager
	funcs  map[string]CallFunc // name => func
}

// New create client
func New() *Pagser {
	p, _ := NewWithConfig(DefaultConfig())
	return p
}

// MustNewWithConfig create client with config
func MustNewWithConfig(cfg Config) *Pagser {
	pagser, err := NewWithConfig(cfg)
	if err != nil {
		panic(err)
	}
	return pagser
}

// NewWithConfig create client with config and error
func NewWithConfig(cfg Config) (*Pagser, error) {
	if cfg.TagerName == "" {
		return nil, errors.New("TagerName must not empty")
	}
	if cfg.FuncSymbol == "" {
		return nil, errors.New("FuncSymbol must not empty")
	}
	if cfg.IgnoreSymbol == "" {
		return nil, errors.New("IgnoreSymbol must not empty")
	}
	return &Pagser{
		config: cfg,
		tagers: make(map[string]*Tager, 0),
		funcs:  sysFuncs,
	}, nil
}

// DefaultConfig the default config
func DefaultConfig() Config {
	return defaultCfg
}
