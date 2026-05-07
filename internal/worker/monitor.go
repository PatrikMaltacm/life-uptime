package worker

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PatrikMaltacm/life-uptime/internal/model"
	"github.com/PatrikMaltacm/life-uptime/internal/repository"
)

type MonitorWorker struct {
	monitorRepo *repository.MonitorRepository
	logRepo     *repository.PingLogRepository
	running     map[string]context.CancelFunc
	mu          sync.Mutex
}

func NewMonitorWorker(mRepo *repository.MonitorRepository, lRepo *repository.PingLogRepository) *MonitorWorker {
	return &MonitorWorker{
		monitorRepo: mRepo,
		logRepo:     lRepo,
		running:     make(map[string]context.CancelFunc),
	}
}

func (w *MonitorWorker) Start(ctx context.Context) {
	log.Println("Worker: Iniciando monitoramento ativo...")

	w.reconcile(ctx)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.reconcile(ctx)
		}
	}
}

func (w *MonitorWorker) reconcile(ctx context.Context) {
	monitors, err := w.monitorRepo.GetAllActive(ctx)
	if err != nil {
		log.Printf("Worker: Erro ao reconciliar monitores: %v", err)
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// inicia monitores novos
	for _, m := range monitors {
		if _, ok := w.running[m.ID]; !ok {
			mCtx, cancel := context.WithCancel(ctx)
			w.running[m.ID] = cancel
			go w.runMonitorCycle(mCtx, m)
			log.Printf("Worker: Monitor %s iniciado", m.URL)
		}
	}

	// para monitores removidos/desativados
	activeIDs := make(map[string]struct{})
	for _, m := range monitors {
		activeIDs[m.ID] = struct{}{}
	}

	for id, cancel := range w.running {
		if _, ok := activeIDs[id]; !ok {
			cancel()
			delete(w.running, id)
			log.Printf("Worker: Monitor %s parado", id)
		}
	}
}

func (w *MonitorWorker) runMonitorCycle(ctx context.Context, m model.MonitorResponse) {
	const defaultInterval = time.Minute

	if m.Interval <= 0 {
		log.Printf("Worker: Monitor %s tem intervalo inválido (%v). Usando padrão de 1m.", m.URL, m.Interval)
		m.Interval = defaultInterval
	}

	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	log.Printf("Worker: Monitorando %s (Intervalo: %v)", m.URL, m.Interval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.pingURL(ctx, m)
		}
	}
}

func (w *MonitorWorker) pingURL(ctx context.Context, m model.MonitorResponse) {
	start := time.Now()

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(m.URL)

	latency := time.Since(start).Milliseconds()

	payload := model.PingLogRequest{
		MonitorID: m.ID,
		Timestamp: time.Now(),
		Latency:   latency,
	}

	if err != nil {
		payload.Error = err.Error()
	} else {
		payload.StatusCode = resp.StatusCode
		resp.Body.Close()
	}

	if err := w.logRepo.Create(ctx, payload); err != nil {
		log.Printf("Worker: erro ao salvar log para %s: %v", m.URL, err)
	}
}
