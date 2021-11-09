package block

import (
	"bytes"
	"fmt"
	"io"

	blocks "github.com/ipfs/go-block-format"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

func NewLinkSystem(svc *Service) ipld.LinkSystem {
	lsys := cidlink.DefaultLinkSystem()
	lsys.TrustedStorage = true
	lsys.StorageReadOpener = blockReadOpener(svc)
	lsys.StorageWriteOpener = blockWriteOpener(svc)
	return lsys
}

func blockReadOpener(svc *Service) ipld.BlockReadOpener {
	return func(lnkCtx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
		asCidLink, ok := lnk.(cidlink.Link)
		if !ok {
			return nil, fmt.Errorf("unsupported link type")
		}
		block, err := svc.GetBlock(lnkCtx.Ctx, asCidLink.Cid)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(block.RawData()), nil
	}
}

func blockWriteOpener(svc *Service) ipld.BlockWriteOpener {
	return func(lnkCtx ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
		buffer := bytes.NewBuffer(nil)
		commit := blockWriteCommitter(svc, buffer)
		return buffer, commit, nil
	}
}

func blockWriteCommitter(svc *Service, buffer *bytes.Buffer) ipld.BlockWriteCommitter {
	return func(lnk ipld.Link) error {
		asCidLink, ok := lnk.(cidlink.Link)
		if !ok {
			return fmt.Errorf("unsupported link type")
		}
		block, err := blocks.NewBlockWithCid(buffer.Bytes(), asCidLink.Cid)
		if err != nil {
			return err
		}
		return svc.AddBlock(block)
	}
}
