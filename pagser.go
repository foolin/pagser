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
	ctxTags  map[string]*tagTokenizer // tag value => tagTokenizer
	ctxFuncs map[string]CallFunc      // name => func
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
	return &Pagser{
		Config:   cfg,
		ctxTags:  make(map[string]*tagTokenizer, 0),
		ctxFuncs: builtinFuncMap,
	}, nil
}
