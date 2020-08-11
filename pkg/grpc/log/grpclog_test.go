package log

import (
	"context"
	conf "github.com/cryptogarageinc/server-common-go/pkg/configuration"
	"github.com/cryptogarageinc/server-common-go/pkg/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLog_Save_IsSaved(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties := `log.format=text
	log.output_stdout=true
	log.level=debug
	`
	config, _ := conf.NewConfigurationFromReader(
		"properties", strings.NewReader(configProperties))
	logConfig := log.Config{}
	config.InitializeComponentConfig(&logConfig)
	l := log.NewLog(&logConfig)
	l.Initialize()
	defer l.Finalize()
	fields := logrus.Fields{"test": "test"}
	entry := l.NewEntry()
	ctx := ctxlogrus.ToContext(context.Background(), entry)
	expected := entry.WithFields(fields)

	// Act
	ctx, _ = Save(ctx, logrus.Fields{"test": "test"})
	actual := ctxlogrus.Extract(ctx)

	// Assert
	assert.Equal(expected, actual)
}
