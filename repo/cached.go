package repo

type cachedRepo struct {
	repo  Repo
	cache map[string]*cacheItem
}

type cacheItem struct {
	pass                string
	lastNotifiedMessage uint64
}

func NewCachedRepo(repo Repo) Repo {
	return &cachedRepo{
		repo:  repo,
		cache: make(map[string]*cacheItem),
	}
}

func (cr *cachedRepo) GetPassword(userName string) (string, error) {
	item, ok := cr.cache[userName]
	if ok && item.pass != "" {
		return item.pass, nil
	}

	pass, err := cr.repo.GetPassword(userName)
	if err != nil {
		return "", err
	}

	// Create the entry if it didn't exist
	if !ok {
		cr.cache[userName] = &cacheItem{}
	}
	cr.cache[userName].pass = pass

	return pass, nil
}

func (cr *cachedRepo) GetLastNotifiedMessage(userName string) (uint64, error) {
	item, ok := cr.cache[userName]
	if ok && item.lastNotifiedMessage != 0 {
		return item.lastNotifiedMessage, nil
	}

	last, err := cr.repo.GetLastNotifiedMessage(userName)
	if err != nil {
		return 0, err
	}

	// Create the entry if it didn't exist
	if !ok {
		cr.cache[userName] = &cacheItem{}
	}
	cr.cache[userName].lastNotifiedMessage = last

	return last, nil
}

func (cr *cachedRepo) SetLastNotifiedMessage(userName string, id uint64) error {
	if err := cr.repo.SetLastNotifiedMessage(userName, id); err != nil {
		return err
	}

	// Create the entry if it didn't exist
	if _, ok := cr.cache[userName]; !ok {
		cr.cache[userName] = &cacheItem{}
	}
	cr.cache[userName].lastNotifiedMessage = id

	return nil
}
