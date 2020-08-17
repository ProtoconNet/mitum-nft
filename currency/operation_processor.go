package currency

import (
	"sync"
	"time"

	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/isaac"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

// BLOCK move ConcurrentOperationsProcessor under isaac.
var maxConcurrentOperations uint = 500

type ConcurrentOperationsProcessor struct {
	size       uint
	sp         *isaac.Statepool
	wk         *util.DistributeWorker
	donechan   chan error
	timeout    time.Duration
	oprLock    sync.RWMutex
	oppHintSet *hint.Hintmap
	oprs       map[hint.Hint]isaac.OperationProcessor
}

func NewConcurrentOperationsProcessor(
	size uint,
	sp *isaac.Statepool,
	timeout time.Duration,
	oppHintSet *hint.Hintmap,
) (*ConcurrentOperationsProcessor, error) {
	if size < 1 {
		return nil, xerrors.Errorf("size must be over 0")
	} else if size > maxConcurrentOperations {
		size = maxConcurrentOperations
	}

	return &ConcurrentOperationsProcessor{
		size:       size,
		sp:         sp,
		timeout:    timeout,
		oppHintSet: oppHintSet,
		oprs:       map[hint.Hint]isaac.OperationProcessor{},
	}, nil
}

func (co *ConcurrentOperationsProcessor) Start() *ConcurrentOperationsProcessor {
	errchan := make(chan error)
	co.wk = util.NewDistributeWorker(co.size, errchan)

	co.donechan = make(chan error, 2)
	go func() {
		<-time.After(co.timeout)
		co.donechan <- xerrors.Errorf("timeout to process")
	}()

	go func() {
		err := co.wk.Run(
			func(i uint, j interface{}) error {
				if j == nil {
					return nil
				} else if op, ok := j.(state.StateProcessor); !ok {
					return xerrors.Errorf("not state.StateProcessor, %T", j)
				} else if opr, err := co.opr(op); err != nil {
					return err
				} else {
					return opr.Process(op)
				}
			},
		)
		if err != nil {
			errchan <- err
		}

		close(errchan)
	}()

	go func() {
		for err := range errchan {
			if err == nil {
				continue
			}

			co.wk.Done(false)
			co.donechan <- err

			return
		}

		co.donechan <- nil
	}()

	return co
}

func (co *ConcurrentOperationsProcessor) Process(po state.StateProcessor) error {
	if co.wk == nil {
		return xerrors.Errorf("not started")
	}

	if !co.wk.NewJob(po) {
		return xerrors.Errorf("already closed")
	}

	return nil
}

func (co *ConcurrentOperationsProcessor) Close() error {
	if co.wk == nil {
		return nil
	}

	co.wk.Done(true)

	return <-co.donechan
}

func (co *ConcurrentOperationsProcessor) opr(op state.StateProcessor) (isaac.OperationProcessor, error) {
	co.oprLock.Lock()
	defer co.oprLock.Unlock()

	var hinter hint.Hinter
	if ht, ok := op.(hint.Hinter); !ok {
		return nil, xerrors.Errorf("not hint.Hinter, %T", op)
	} else {
		hinter = ht
	}

	if opr, found := co.oprs[hinter.Hint()]; found {
		return opr, nil
	}

	var opr isaac.OperationProcessor
	if hinter, found := co.oppHintSet.Get(hinter); !found {
		opr = defaultOperationProcessor{}
	} else {
		opr = hinter.(isaac.OperationProcessor)
	}

	opr = opr.New(co.sp)
	co.oprs[hinter.Hint()] = opr

	return opr, nil
}

type defaultOperationProcessor struct {
	pool *isaac.Statepool
}

func (opp defaultOperationProcessor) New(pool *isaac.Statepool) isaac.OperationProcessor {
	return &defaultOperationProcessor{
		pool: pool,
	}
}

func (opp defaultOperationProcessor) Process(op state.StateProcessor) error {
	return op.Process(opp.pool.Get, opp.pool.Set)
}

type PreProcessor interface {
	PreProcess(
		getState func(key string) (state.StateUpdater, bool, error),
		setState func(valuehash.Hash, ...state.StateUpdater) error,
	) error
}

type OperationProcessor struct {
	sync.RWMutex
	pool            *isaac.Statepool
	amountPool      map[string]*AmountState
	processedSender map[string]struct{}
}

func (opr *OperationProcessor) New(pool *isaac.Statepool) isaac.OperationProcessor {
	return &OperationProcessor{
		pool:            pool,
		amountPool:      map[string]*AmountState{},
		processedSender: map[string]struct{}{},
	}
}

func (opr *OperationProcessor) getState(key string) (state.StateUpdater, bool, error) {
	opr.Lock()
	defer opr.Unlock()

	if ast, found := opr.amountPool[key]; found {
		return ast, true, nil
	} else if st, found, err := opr.pool.Get(key); err != nil {
		return nil, false, err
	} else {
		ast := NewAmountState(st)
		opr.amountPool[key] = ast

		return ast, found, nil
	}
}

func (opr *OperationProcessor) Process(op state.StateProcessor) error {
	var sp state.StateProcessor
	var sender string
	var get func(string) (state.StateUpdater, bool, error)

	switch t := op.(type) {
	case Transfer:
		get = opr.getState
		sp = &TransferProcessor{Transfer: t}
		sender = t.Fact().(TransferFact).Sender().String()
	case CreateAccount:
		get = opr.getState
		sp = &CreateAccountProcessor{CreateAccount: t}
		sender = t.Fact().(CreateAccountFact).Sender().String()
	default:
		return sp.Process(opr.pool.Get, opr.pool.Set)
	}

	if func() bool {
		opr.RLock()
		defer opr.RUnlock()

		_, found := opr.processedSender[sender]

		return found
	}() {
		return state.IgnoreOperationProcessingError.Errorf("violates only one sender in proposal")
	}

	if pr, ok := sp.(PreProcessor); ok {
		if err := pr.PreProcess(get, opr.pool.Set); err != nil {
			return err
		}
	}

	if err := sp.Process(get, opr.pool.Set); err != nil {
		return err
	} else {
		opr.Lock()
		defer opr.Unlock()

		opr.processedSender[sender] = struct{}{}

		return nil
	}
}
