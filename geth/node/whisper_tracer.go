package node

import (
	"strconv"

	whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	envelopeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "envelope_counter",
			Help: "Envelopes counter",
		},
		[]string{"topic", "source", "is_new", "peer"},
	)
	envelopeVolume = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "envelope_volume",
			Help: "Volume of received envelopes",
		},
		[]string{"topic", "source", "is_new", "peer"},
	)
)

func init() {
	prometheus.MustRegister(envelopeCounter)
	prometheus.MustRegister(envelopeVolume)
}

// WhisperEnvelopeTracer traces incoming envelopes.
type WhisperEnvelopeTracer struct{}

// Trace is called for every incoming envelope.
func (t *WhisperEnvelopeTracer) Trace(envelope *whisper.EnvelopeMeta) {
	envelopeCounter.WithLabelValues(
		envelope.Topic.String(),
		envelope.SourceString(),
		strconv.FormatBool(envelope.IsNew),
		envelope.Peer,
	).Inc()
	envelopeVolume.WithLabelValues(
		envelope.Topic.String(),
		envelope.SourceString(),
		strconv.FormatBool(envelope.IsNew),
		envelope.Peer,
	).Add(float64(envelope.Size))
}
