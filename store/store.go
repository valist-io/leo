package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

type Store struct {
	bs blockstore.Blockstore
	ex exchange.Interface
}

func NewStore(bs blockstore.Blockstore, ex exchange.Interface) *Store {
	return &Store{bs, ex}
}

func (s *Store) AddBlock(block blocks.Block) error {
	if has, err := s.bs.Has(block.Cid()); has || err != nil {
		return err
	}

	if err := s.bs.Put(block); err != nil {
		return err
	}

	return s.ex.HasBlock(block)
}

func (s *Store) GetBlock(ctx context.Context, id cid.Cid) (blocks.Block, error) {
	block, err := s.bs.Get(id)
	if err == nil {
		return block, nil
	}

	if err == blockstore.ErrNotFound {
		return s.ex.GetBlock(ctx, id)
	}

	return nil, err
}

func (s *Store) LinkSystem() ipld.LinkSystem {
	lsys := cidlink.DefaultLinkSystem()
	lsys.TrustedStorage = true
	lsys.StorageReadOpener = s.blockReadOpener
	lsys.StorageWriteOpener = s.blockWriteOpener
	return lsys
}

func (s *Store) blockWriteOpener(lnkCtx ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
	bwc := &blockWriteCommitter{store: s}
	return bwc, bwc.Commit, nil
}

func (s *Store) blockReadOpener(lnkCtx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
	asCidLink, ok := lnk.(cidlink.Link)
	if !ok {
		return nil, fmt.Errorf("unsupported link type")
	}

	block, err := s.GetBlock(lnkCtx.Ctx, asCidLink.Cid)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(block.RawData()), nil
}

type blockWriteCommitter struct {
	store  *Store
	buffer bytes.Buffer
}

func (bwc *blockWriteCommitter) Write(data []byte) (int, error) {
	return bwc.buffer.Write(data)
}

func (bwc *blockWriteCommitter) Commit(lnk ipld.Link) error {
	asCidLink, ok := lnk.(cidlink.Link)
	if !ok {
		return fmt.Errorf("unsupported link type")
	}

	block, err := blocks.NewBlockWithCid(bwc.buffer.Bytes(), asCidLink.Cid)
	if err != nil {
		return err
	}

	return bwc.store.AddBlock(block)
}
