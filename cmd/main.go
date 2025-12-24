package main

import (
	wasmgoap "accounter/adapters/ui/wasm-goap"
	"accounter/config"
	"accounter/internal/app"
	"accounter/pkg/logger"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Init graceful shutdown context
	ctx, cancel := config.InitGracefulShutdownCtx()

	// Init application config
	cfg := config.InitConfig()

	// Create logger
	logger := logger.NewLogger(cfg.DebugMode, cfg.AppMode, "logs")

	// Create frontend application instance
	frontApp := wasmgoap.NewApp(cfg, logger)

	// Create AppContext instance and init it
	backendApp := app.NewAppContext(ctx, cfg, logger).Init(ctx)

	// Register background tasks
	backendApp.RegisterTask("Frontend HTTP server", &frontApp)

	// Defer canceling and shutdown application
	defer func() {
		cancel()
		backendApp.Shutdown()
	}()

	// Run backend application
	backendApp.Run(ctx)
}

var payments = []string{
	"11790;22/12/2025",
	"ОПЛАТА ПО СЧЕТУ №11785 ОТ 18.12.2025 Г. ДОГОВОР № 668 ОТ 14.11.2025 Г. СУММА 13360-00 БЕЗ НАЛОГА (НДС) ",
	"ОПЛАТА ПО СЧЕТУ 11643 ОТ 01.12.2025 ЗА АБОНЕНТСКОГО ОБСЛ. ЗА ДЕКАБРЬ 2025Г. СУММА 85050-00 БЕЗ НАЛОГА (НДС)",
	"Оплата по счету  11747 от 16.12.2025  за АО (использование ПО, услуги мониторинга)  Сумма 119350-00 Без налога (НДС)",
}

func parsePayment(value string) (number int, date time.Time, err error) {
	value = strings.TrimSpace(strings.ToLower(value))
	tz := time.FixedZone("Custom", 180*60)

	regPhysic := regexp.MustCompile(`\d+;[1-31]+/[1-12]+/\d{4,}$`)
	regJur := regexp.MustCompile(`оплата по счету\s+[№]?\d+ от \d{1,2}.\d{1,2}.\d{4,}`)

	switch {
	case regPhysic.MatchString(value):
		return parseNumberAndDate("02/01/2006", tz, strings.Split(value, ";")...)

	case regJur.MatchString(value):
		item := regJur.FindString(value)
		item = strings.NewReplacer("оплата по счету ", "", " от", "", "№", "").Replace(item)
		item = strings.TrimSpace(item)

		return parseNumberAndDate("02.01.2006", tz, strings.Split(item, " ")...)

	default:
		err = errors.New("bad payment format")
	}

	return
}

func parseNumberAndDate(layout string, tz *time.Location, items ...string) (number int, date time.Time, err error) {
	switch len(items) {
	case 2:
		date, err = time.ParseInLocation(layout, items[1], tz)
		fallthrough

	case 1:
		number, err = strconv.Atoi(items[0])

	default:
		err = fmt.Errorf("invalid payment data items length: %d", len(items))
	}

	return
}
