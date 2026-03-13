package sas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BloomFilterFile = "bloomfilter.json"
	LRUFile         = "lrucache.json"
	BlackListFile   = "blacklist.json"
	StatsFile       = "stats.json"
)

type Keeper struct {
	coinKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	dataDir    string
}

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, dataDir string) Keeper {
	k := Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
		dataDir:    dataDir,
	}
	k.loadBloomFilter()
	k.loadLRUCache()
	k.loadBlackList()
	k.loadStats()
	return k
}

func (k Keeper) GetLAddress(ctx sdk.Context, sUrl string) LAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(sUrl))
	if bz == nil {
		return LAddress{}
	}
	var lAddress LAddress
	k.cdc.MustUnmarshalBinaryBare(bz, &lAddress)
	return lAddress
}

func (k Keeper) StoreLAddress(ctx sdk.Context, sUrl string, owner sdk.AccAddress, price sdk.Coins, duration time.Duration) {
	lAddress := NewLAddress()
	lAddress.Price = price
	lAddress.Owner = owner
	lAddress.ExpirationTime = time.Now().Add(duration)
	lAddress.ClickCount = 0
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(sUrl), k.cdc.MustMarshalBinaryBare(lAddress))
	k.AddToBloomFilter(sUrl)
}

func (k Keeper) isSUrlExist(ctx sdk.Context, sUrl string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(sUrl))
}

func (k Keeper) SetLAddress(ctx sdk.Context, sUrl string, lAddress LAddress) {
	if lAddress.Owner.Empty() {
		return
	}
	lAddress.UpdatedAt = time.Now()
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(sUrl), k.cdc.MustMarshalBinaryBare(lAddress))
}

func (k Keeper) ResolveLUrl(ctx sdk.Context, sUrl string) string {
	lAddress := k.GetLAddress(ctx, sUrl)
	if lAddress.LUrl == "" {
		return ""
	}
	k.IncrementClickCount(ctx, sUrl)
	return lAddress.LUrl
}

func (k Keeper) SetLUrl(ctx sdk.Context, sUrl string, lUrl string) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.LUrl = lUrl
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) HasOwner(ctx sdk.Context, sUrl string) bool {
	return !k.GetLAddress(ctx, sUrl).Owner.Empty()
}

func (k Keeper) GetOwner(ctx sdk.Context, sUrl string) sdk.AccAddress {
	return k.GetLAddress(ctx, sUrl).Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, sUrl string, owner sdk.AccAddress) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.Owner = owner
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) GetPrice(ctx sdk.Context, sUrl string) sdk.Coins {
	return k.GetLAddress(ctx, sUrl).Price
}

func (k Keeper) SetPrice(ctx sdk.Context, sUrl string, price sdk.Coins) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.Price = price
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) GetSUrlsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k Keeper) SetSell(ctx sdk.Context, sUrl string, isSell bool) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.IsSell = isSell
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) SetNoSell(ctx sdk.Context, sUrl string) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.IsSell = false
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
}

func (k Keeper) GetSell(ctx sdk.Context, sUrl string) bool {
	return k.GetLAddress(ctx, sUrl).IsSell
}

func (k Keeper) IsExpired(ctx sdk.Context, sUrl string) bool {
	lAddress := k.GetLAddress(ctx, sUrl)
	return time.Now().After(lAddress.ExpirationTime)
}

func (k Keeper) IsInGracePeriod(ctx sdk.Context, sUrl string) bool {
	lAddress := k.GetLAddress(ctx, sUrl)
	return time.Now().Before(lAddress.ExpirationTime.Add(GracePeriod))
}

func (k Keeper) GetExpirationTime(ctx sdk.Context, sUrl string) time.Time {
	return k.GetLAddress(ctx, sUrl).ExpirationTime
}

func (k Keeper) Renew(ctx sdk.Context, sUrl string, duration time.Duration) sdk.Error {
	lAddress := k.GetLAddress(ctx, sUrl)
	if time.Now().After(lAddress.ExpirationTime) {
		lAddress.ExpirationTime = time.Now().Add(duration)
	} else {
		lAddress.ExpirationTime = lAddress.ExpirationTime.Add(duration)
	}
	lAddress.UpdatedAt = time.Now()
	k.SetLAddress(ctx, sUrl, lAddress)
	return nil
}

func (k Keeper) DeleteLAddress(ctx sdk.Context, sUrl string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(sUrl))
	k.RemoveFromBloomFilter(sUrl)
}

func (k Keeper) CleanExpiredAddresses(ctx sdk.Context) (count int) {
	iterator := k.GetSUrlsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		sUrl := string(iterator.Key())
		lAddress := k.GetLAddress(ctx, sUrl)
		if time.Now().After(lAddress.ExpirationTime.Add(GracePeriod)) {
			k.DeleteLAddress(ctx, sUrl)
			count++
		}
	}
	return count
}

func (k Keeper) GetClickCount(ctx sdk.Context, sUrl string) uint64 {
	return k.GetLAddress(ctx, sUrl).ClickCount
}

func (k Keeper) IncrementClickCount(ctx sdk.Context, sUrl string) {
	lAddress := k.GetLAddress(ctx, sUrl)
	lAddress.ClickCount++
	k.SetLAddress(ctx, sUrl, lAddress)
	k.incrementTotalClicks(ctx)
}

func (k Keeper) GetAllExpiredWithGracePeriod(ctx sdk.Context) []ExpiredSUrl {
	var expired []ExpiredSUrl
	iterator := k.GetSUrlsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		sUrl := string(iterator.Key())
		lAddress := k.GetLAddress(ctx, sUrl)
		if time.Now().After(lAddress.ExpirationTime.Add(GracePeriod)) {
			expired = append(expired, ExpiredSUrl{
				SUrl:        sUrl,
				ExpiredAt:   lAddress.ExpirationTime,
				OriginalExp: lAddress.ExpirationTime,
			})
		}
	}
	return expired
}

func (k Keeper) IsBlackListed(url string) bool {
	bl := k.getBlackListFromStore()
	if bl.URLs[url] {
		return true
	}
	parsedURL := url
	if !strings.HasPrefix(parsedURL, "http") {
		parsedURL = "https://" + url
	}
	for domain := range bl.Domains {
		if strings.Contains(parsedURL, domain) {
			return true
		}
	}
	return false
}

func (k Keeper) getBlackListFromStore() BlackList {
	return BlackList{
		URLs:    make(map[string]bool),
		Domains: make(map[string]bool),
	}
}

func (k Keeper) GetBlackList(ctx sdk.Context) BlackList {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("blacklist"))
	if bz == nil {
		return NewBlackList()
	}
	var bl BlackList
	k.cdc.MustUnmarshalJSON(bz, &bl)
	return bl
}

func (k Keeper) SetBlackList(ctx sdk.Context, bl BlackList) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte("blacklist"), k.cdc.MustMarshalJSON(bl))
}

func (k Keeper) AddToBlackListURL(ctx sdk.Context, url string) {
	bl := k.GetBlackList(ctx)
	bl.URLs[url] = true
	bl.UpdatedAt = time.Now()
	k.SetBlackList(ctx, bl)
	k.saveBlackList()
}

func (k Keeper) AddToBlackListDomain(ctx sdk.Context, domain string) {
	bl := k.GetBlackList(ctx)
	bl.Domains[domain] = true
	bl.UpdatedAt = time.Now()
	k.SetBlackList(ctx, bl)
	k.saveBlackList()
}

func (k Keeper) RemoveFromBlackListURL(ctx sdk.Context, url string) {
	bl := k.GetBlackList(ctx)
	delete(bl.URLs, url)
	bl.UpdatedAt = time.Now()
	k.SetBlackList(ctx, bl)
	k.saveBlackList()
}

func (k Keeper) GetStats(ctx sdk.Context) Stats {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("stats"))
	if bz == nil {
		return NewStats()
	}
	var stats Stats
	k.cdc.MustUnmarshalJSON(bz, &stats)
	return stats
}

func (k Keeper) SetStats(ctx sdk.Context, stats Stats) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte("stats"), k.cdc.MustMarshalJSON(stats))
}

func (k Keeper) incrementTotalClicks(ctx sdk.Context) {
	stats := k.GetStats(ctx)
	stats.TotalClicks++
	today := time.Now().Format("2006-01-02")
	stats.DailyClicks[today]++
	k.updateTopShortURLs(ctx)
	k.SetStats(ctx, stats)
	k.saveStatsToFile()
}

func (k Keeper) updateTopShortURLs(ctx sdk.Context) {
	stats := k.GetStats(ctx)
	var entries []StatEntry
	iterator := k.GetSUrlsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		sUrl := string(iterator.Key())
		lAddress := k.GetLAddress(ctx, sUrl)
		entries = append(entries, StatEntry{Key: sUrl, Value: lAddress.ClickCount})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	if len(entries) > 10 {
		entries = entries[:10]
	}
	stats.TopShortURLs = entries
	k.SetStats(ctx, stats)
}

func (k Keeper) GetTopStats(ctx sdk.Context) Stats {
	return k.GetStats(ctx)
}

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	count := k.CleanExpiredAddresses(ctx)
	if count > 0 {
		fmt.Printf("Cleaned %d expired addresses\n", count)
	}
	k.saveBloomFilter()
	k.saveLRUCache()
}

func (k Keeper) saveBloomFilter() {
	if k.dataDir == "" {
		return
	}
	data, _ := json.Marshal(GlobalBloomFilter)
	ioutil.WriteFile(filepath.Join(k.dataDir, BloomFilterFile), data, 0644)
}

func (k Keeper) loadBloomFilter() {
	if k.dataDir == "" {
		return
	}
	file := filepath.Join(k.dataDir, BloomFilterFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	GlobalBloomFilter = NewBloomFilter(1024*1024, 16)
	json.Unmarshal(data, GlobalBloomFilter)
}

func (k Keeper) saveLRUCache() {
	if k.dataDir == "" || LruCache == nil {
		return
	}
	data, _ := json.Marshal(LruCache)
	ioutil.WriteFile(filepath.Join(k.dataDir, LRUFile), data, 0644)
}

func (k Keeper) loadLRUCache() {
	if k.dataDir == "" {
		return
	}
	file := filepath.Join(k.dataDir, LRUFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	LruCache = New(10000)
	json.Unmarshal(data, LruCache)
}

func (k Keeper) saveBlackList() {
	if k.dataDir == "" {
		return
	}
	bl := BlackList{
		URLs:      make(map[string]bool),
		Domains:   make(map[string]bool),
		UpdatedAt: time.Now(),
	}
	data, _ := json.Marshal(bl)
	ioutil.WriteFile(filepath.Join(k.dataDir, BlackListFile), data, 0644)
}

func (k Keeper) loadBlackList() {
	if k.dataDir == "" {
		return
	}
	file := filepath.Join(k.dataDir, BlackListFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	bl := NewBlackList()
	json.Unmarshal(data, &bl)
	GlobalBlackList = &bl
}

func (k Keeper) saveStatsToFile() {
	if k.dataDir == "" {
		return
	}
	stats := Stats{
		TotalClicks:  0,
		DailyClicks:  make(map[string]uint64),
		TopShortURLs: make([]StatEntry, 0),
		TopOwners:    make([]StatEntry, 0),
	}
	data, _ := json.Marshal(stats)
	ioutil.WriteFile(filepath.Join(k.dataDir, StatsFile), data, 0644)
}

func (k Keeper) loadStats() {
	if k.dataDir == "" {
		return
	}
	file := filepath.Join(k.dataDir, StatsFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	stats := NewStats()
	json.Unmarshal(data, &stats)
	GlobalStats = &stats
}

var GlobalBlackList *BlackList
var GlobalStats *Stats

func init() {
	GlobalBlackList = new(BlackList)
	GlobalBlackList.URLs = make(map[string]bool)
	GlobalBlackList.Domains = make(map[string]bool)
	GlobalStats = new(Stats)
	GlobalStats.DailyClicks = make(map[string]uint64)
	GlobalStats.TopShortURLs = make([]StatEntry, 0)
	GlobalStats.TopOwners = make([]StatEntry, 0)
}
