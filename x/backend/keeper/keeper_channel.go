package keeper

import (
	"fmt"
	"time"

	"github.com/okex/okchain/x/backend/config"
	"github.com/okex/okchain/x/backend/orm"
	"github.com/okex/okchain/x/backend/types"

	"github.com/tendermint/tendermint/libs/log"
)

func pushAllKline1m(klines map[string][]types.KlineM1, keeper Keeper) {
	keeper.Logger.Debug("pushAllKline1m_1", "klines", klines)
	if klines != nil && len(klines) > 0 {
		for _, klineArr := range klines {
			if klineArr == nil {
				continue
			}

			for _, k := range klineArr {
				keeper.Logger.Debug("pushAllKline1m_2", "kline", &k)
				keeper.pushWSItem(&k)
			}
		}
	}
}

func generateKline1M(stop chan struct{}, conf *config.Config, o *orm.ORM, log *log.Logger, keeper Keeper) {
	o.Debug("[backend] generateKline1M go routine started")
	defer types.PrintStackIfPanic()

	startTS, endTS := int64(0), time.Now().Unix()-60
	time.Sleep(3 * time.Second)
	if o.GetMaxBlockTimestamp() > 0 {
		endTS = o.GetMaxBlockTimestamp()
	}

	//ds := DealDataSource{orm: orm}
	ds := orm.MergeResultDataSource{Orm: o}
	anchorEndTS, _, newKline1s, err := o.CreateKline1min(startTS, endTS, &ds)
	if err != nil {
		(*log).Debug(fmt.Sprintf("[backend] error: %+v \n", err))
	}

	pushAllKline1m(newKline1s, keeper)

	waitInSecond := int(60+types.Kline1GoRoutineWaitInSecond-time.Now().Second()) % 60
	timer := time.NewTimer(time.Duration(waitInSecond * int(time.Second)))
	interval := time.Second * 60
	ticker := time.NewTicker(interval)

	go CleanUpKlines(stop, o, conf)
	var klineNotifyChans *map[int]chan struct{}
	work := func() {
		if o.GetMaxBlockTimestamp() == 0 {
			return
		}

		crrtBlkTS := o.GetMaxBlockTimestamp()
		(*log).Debug(fmt.Sprintf("[backend] entering generateKline1M [%d, %d) [%s, %s)",
			anchorEndTS, crrtBlkTS, types.TimeString(anchorEndTS), types.TimeString(crrtBlkTS)))

		anchorStart, _, newKline1s, err := o.CreateKline1min(anchorEndTS, crrtBlkTS, &ds)
		(*log).Debug(fmt.Sprintf("[backend] generateKline1M's actually merge period [%s, %s)",
			types.TimeString(anchorEndTS), types.TimeString(anchorStart)))
		if err != nil {
			(*log).Debug(fmt.Sprintf("[backend] generateKline1M error: %s", err.Error()))

		} else {
			pushAllKline1m(newKline1s, keeper)
			anchorEndTS = anchorStart
			if klineNotifyChans != nil {
				for _, ch := range *klineNotifyChans {
					ch <- struct{}{}
				}
			}

		}
	}

	work()

	klineNotifyChans = generateSyncKlineMXChans()
	for freq, ntfCh := range *klineNotifyChans {
		go generateKlinesMX(ntfCh, stop, freq, o, keeper)
	}

	for {
		select {
		case <-timer.C:
			work()
			ticker = time.NewTicker(interval)
		case <-ticker.C:
			work()
		case <-stop:
			break

		}
	}
}

func generateSyncKlineMXChans() *map[int]chan struct{} {
	notifyChans := map[int]chan struct{}{}
	klineMap := types.GetAllKlineMap()

	for freq := range klineMap {
		if freq > 60 {
			notifyCh := make(chan struct{}, 1)
			notifyChans[freq] = notifyCh
		}
	}

	return &notifyChans
}


func pushAllKlineXm(klines map[string][]interface{}, keeper Keeper) {
	if klines != nil && len(klines) >0 {
		for _, klineArr := range klines {
			if klineArr == nil {
				continue
			}

			for _, k := range klineArr {
				baseLine := k.(types.IWebsocket)
				keeper.pushWSItem(baseLine)
			}
		}
	}
}

func generateKlinesMX(notifyChan chan struct{}, stop chan struct{}, refreshInterval int, o *orm.ORM, keeper Keeper) {
	o.Debug(fmt.Sprintf("[backend] generateKlineMX-#%d# go routine started", refreshInterval))

	destKName := types.GetKlineTableNameByFreq(refreshInterval)
	destK, err := types.NewKlineFactory(destKName, nil)
	if err != nil {
		o.Error(fmt.Sprintf("[backend] NewKlineFactory error: %s", err.Error()))
	}

	destIKline := destK.(types.IKline)

	startTS, endTS := int64(0), time.Now().Unix()-int64(destIKline.GetFreqInSecond())
	anchorEndTS, _, newKlines, err := o.MergeKlineM1(startTS, endTS, destIKline)
	if err != nil {
		o.Error(fmt.Sprintf("[backend] MergeKlineM1 error: %s", err.Error()))
	}
	pushAllKlineXm(newKlines, keeper)

	//waitInSecond := int(60+KlineX_GOROUTINE_WAIT_IN_SECOND-time.Now().Second()) % 60
	crrTS := time.Now().Unix()
	waitInSecond := int64(destIKline.GetFreqInSecond()) - (crrTS - destIKline.GetAnchorTimeTS(crrTS)) + types.KlinexGoRoutineWaitInSecond + 60
	timer := time.NewTimer(time.Duration(int(waitInSecond) * int(time.Second)))
	interval := time.Duration(destIKline.GetFreqInSecond() * int(time.Second))
	o.Debug(fmt.Sprintf("[backend] duaration: %+v(%d s) IKline: %+v ", interval, destIKline.GetFreqInSecond(), destIKline))
	ticker := time.NewTicker(interval)

	work := func() {
		if o.GetMaxBlockTimestamp() == 0 {
			return
		}

		crrtTS := o.GetMaxBlockTimestamp()

		o.Debug(fmt.Sprintf("[backend] entering generateKlinesMX-#%d# [%d, %d)[%s, %s)",
			destIKline.GetFreqInSecond(), anchorEndTS, crrtTS, types.TimeString(anchorEndTS), types.TimeString(crrtTS)))

		anchorStart, _, newKlines, err := o.MergeKlineM1(anchorEndTS, crrtTS, destIKline)
		if err != nil {
			o.Error(fmt.Sprintf("[backend] error: %s", err.Error()))

		} else {
			anchorEndTS = anchorStart
			pushAllKlineXm(newKlines, keeper)
		}
	}

	for {
		select {
		case <-notifyChan:
			time.Sleep(time.Second)
			if anchorEndTS > 0 && time.Now().Unix() < anchorEndTS+int64(destIKline.GetFreqInSecond()) {
				break
			} else {
				work()
				ticker = time.NewTicker(interval)
			}

		case <-ticker.C:
			work()
		case <-timer.C:
			work()
		case <-stop:
			break
		}
	}
}

// nolint
func CleanUpKlines(stop chan struct{}, o *orm.ORM, conf *config.Config) {
	o.Debug(fmt.Sprintf("[backend] cleanUpKlines go routine started. MaintainConf: %+v", *conf))
	interval := time.Duration(60 * int(time.Second))
	ticker := time.NewTicker(time.Duration(int(60-time.Now().Second()) * int(time.Second)))

	work := func() {
		now := time.Now()
		strNow := now.Format("15:04:05")
		if strNow == conf.CleanUpsTime {

			m := types.GetAllKlineMap()
			for _, ktype := range m {
				expiredDays := conf.CleanUpsKeptDays[ktype]
				if expiredDays != 0 {
					o.Debug(fmt.Sprintf("[backend] entering cleanUpKlines, "+
						"fired time: %s(currentTS: %d), kline type: %s", conf.CleanUpsTime, now.Unix(), ktype))
					//anchorTS := now.Add(-time.Duration(int(time.Second) * 1440 * expiredDays)).Unix()
					anchorTS := now.Add(-time.Duration(int(time.Second) * types.SecondsInADay * expiredDays)).Unix()
					kline, err := types.NewKlineFactory(ktype, nil)
					if err != nil {
						o.Debug("failed to NewKlineFactory becoz of : " + err.Error())
						break
					}
					if err = o.DeleteKlineBefore(anchorTS, kline); err != nil {
						o.Error("failed to DeleteKlineBefore because " + err.Error())
					}
				}
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			work()
			ticker = time.NewTicker(interval)

		case <-stop:
			break

		}
	}
}
