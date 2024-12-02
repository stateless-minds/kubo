package node

import (
	blockstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/go-datastore"
	config "github.com/ipfs/kubo/config"
	"go.uber.org/fx"

	"github.com/ipfs/boxo/filestore"
	"github.com/ipfs/kubo/core/node/helpers"
	"github.com/ipfs/kubo/repo"
	"github.com/ipfs/kubo/thirdparty/verifbs"
)

// RepoConfig loads configuration from the repo
func RepoConfig(repo repo.Repo) (*config.Config, error) {
	cfg, err := repo.Config()
	return cfg, err
}

// Datastore provides the datastore
func Datastore(repo repo.Repo) datastore.Datastore {
	return repo.Datastore()
}

// BaseBlocks is the lower level blockstore without GC or Filestore layers
type BaseBlocks blockstore.Blockstore

// BaseBlockstoreCtor creates cached blockstore backed by the provided datastore
func BaseBlockstoreCtor(cacheOpts blockstore.CacheOpts, hashOnRead bool) func(mctx helpers.MetricsCtx, repo repo.Repo, lc fx.Lifecycle) (bs BaseBlocks, err error) {
	return func(mctx helpers.MetricsCtx, repo repo.Repo, lc fx.Lifecycle) (bs BaseBlocks, err error) {
		// hash security
		bs = blockstore.NewBlockstore(repo.Datastore())
		bs = &verifbs.VerifBS{Blockstore: bs}
		bs, err = blockstore.CachedBlockstore(helpers.LifecycleCtx(mctx, lc), bs, cacheOpts)
		if err != nil {
			return nil, err
		}

		bs = blockstore.NewIdStore(bs)

		if hashOnRead { // TODO: review: this is how it was done originally, is there a reason we can't just pass this directly?
			bs.HashOnRead(true)
		}

		return
	}
}

// GcBlockstoreCtor wraps the base blockstore with GC and Filestore layers
func GcBlockstoreCtor(bb BaseBlocks) (gclocker blockstore.GCLocker, gcbs blockstore.GCBlockstore, bs blockstore.Blockstore) {
	gclocker = blockstore.NewGCLocker()
	gcbs = blockstore.NewGCBlockstore(bb, gclocker)

	bs = gcbs
	return
}

// FilestoreBlockstoreCtor wraps GcBlockstore and adds Filestore support
func FilestoreBlockstoreCtor(repo repo.Repo, bb BaseBlocks) (gclocker blockstore.GCLocker, gcbs blockstore.GCBlockstore, bs blockstore.Blockstore, fstore *filestore.Filestore) {
	gclocker = blockstore.NewGCLocker()

	// hash security
	fstore = filestore.NewFilestore(bb, repo.FileManager())
	gcbs = blockstore.NewGCBlockstore(fstore, gclocker)
	gcbs = &verifbs.VerifBSGC{GCBlockstore: gcbs}

	bs = gcbs
	return
}
