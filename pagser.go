// Copyright 2020 Foolin

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pagser

import (
	"errors"
	"sync"
)

// Pagser the page parser
type Pagser struct {
	Config Config
	//mapTags  map[string]*tagTokenizer // tag value => tagTokenizer
	mapTags sync.Map //map[string]*tagTokenizer
	//mapFuncs map[string]CallFunc      // name => func
	mapFuncs sync.Map //map[string]CallFunc
}

// New create pagser client
func New() *Pagser {
	p, _ := NewWithConfig(DefaultConfig())
	return p
}

// NewWithConfig create pagser client with Config and error
func NewWithConfig(cfg Config) (*Pagser, error) {
	if cfg.TagName == "" {
		return nil, errors.New("tag name must not empty")
	}
	if cfg.FuncSymbol == "" {
		return nil, errors.New("FuncSymbol must not empty")
	}
	p := Pagser{
		Config: cfg,
		//mapTags:  make(map[string]*tagTokenizer, 0),
		//mapFuncs: builtinFuncs,
	}
	for k, v := range builtinFuncs {
		p.mapFuncs.Store(k, v)
	}
	return &p, nil
}
