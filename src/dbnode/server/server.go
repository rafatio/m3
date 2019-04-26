	"github.com/m3db/m3/src/dbnode/encoding/proto"
	xconfig "github.com/m3db/m3/src/x/config"
	"github.com/m3db/m3/src/x/context"
	"github.com/m3db/m3/src/x/ident"
	"github.com/m3db/m3/src/x/instrument"
	"github.com/m3db/m3/src/x/pool"
	xsync "github.com/m3db/m3/src/x/sync"
	"github.com/jhump/protoreflect/desc"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	defaultServiceName               = "m3dbnode"
	defer logger.Sync()
		logger.Fatal("could not parse new file mode", zap.Error(err))
		logger.Fatal("could not parse new directory mode", zap.Error(err))
		logger.Fatal("could not acquire lock", zap.String("path", lockPath), zap.Error(err))
	var schema *desc.MessageDescriptor
	if cfg.Proto != nil {
		logger.Info("Probuf data mode enabled")
		schema, err = proto.ParseProtoSchema(
			cfg.Proto.SchemaFilePath, cfg.Proto.MessageName)
		if err != nil {
			logger.Fatal("error parsing protobuffer schema", zap.Error(err))
		}
	}

		logger.Fatal("could not connect to metrics", zap.Error(err))
		logger.Fatal("could not resolve local host ID", zap.Error(err))
	}

	var (
		tracer      opentracing.Tracer
		traceCloser io.Closer
	)

	if cfg.Tracing == nil {
		tracer = opentracing.NoopTracer{}
		logger.Info("tracing disabled; set `tracing.backend` to enable")
	} else {
		// setup tracer
		serviceName := cfg.Tracing.ServiceName
		if serviceName == "" {
			serviceName = defaultServiceName
		}
		tracer, traceCloser, err = cfg.Tracing.NewTracer(serviceName, scope.SubScope("jaeger"), logger)
		if err != nil {
			tracer = opentracing.NoopTracer{}
			logger.Warn("could not initialize tracing; using no-op tracer instead",
				zap.String("service", serviceName), zap.Error(err))
		} else {
			defer traceCloser.Close()
			logger.Info("tracing enabled", zap.String("service", serviceName))
		}
				logger.Fatal("unable to create etcd clusters", zap.Error(err))
			logger.Info("using seed nodes etcd cluster",
				zap.String("zone", zone), zap.Strings("endpoints", endpoints))
		logger.Info("resolving seed node configuration",
			zap.String("hostID", hostID), zap.Strings("seedNodeHostIDs", seedNodeHostIDs),
		)
				logger.Fatal("unable to create etcd config", zap.Error(err))
				logger.Fatal("could not start embedded etcd", zap.Error(err))
		SetMetricsSamplingRate(cfg.Metrics.SampleRate()).
		SetTracer(tracer)
	opentracing.SetGlobalTracer(tracer)

		logger.Warn("max index query IDs concurrency was not set, falling back to default value")
		logger.Fatal("unable to start build reporter", zap.Error(err))
		logger.Fatal("could not construct postings list cache", zap.Error(err))
		logger.Fatal("could not set initial runtime options", zap.Error(err))
			logger.Fatal("could not determine if host supports HugeTLB", zap.Error(err))
			logger.Warn("host doesn't support HugeTLB, proceeding without it")
		logger.Fatal("unknown commit log queue size type",
			zap.Any("type", cfg.CommitLog.Queue.CalculationType))
			logger.Fatal("unknown commit log queue channel size type",
				zap.Any("type", cfg.CommitLog.Queue.CalculationType))
	opts = withEncodingAndPoolingOptions(cfg, logger, schema, opts, cfg.PoolingPolicy)
		logger.Fatal("could not create persist manager", zap.Error(err))
			logger.Fatal("could not initialize dynamic config", zap.Error(err))
			logger.Fatal("could not initialize static config", zap.Error(err))
		logger.Fatal("could not initialize m3db topology", zap.Error(err))
		},
		func(opts client.AdminOptions) client.AdminOptions {
			if cfg.Proto != nil {
				return opts.SetEncodingProto(
					schema,
					encoding.NewOptions(),
				).(client.AdminOptions)
			}
			return opts
		},
	)
		logger.Fatal("could not create m3db client", zap.Error(err))
		logger.Fatal("could not create bootstrap process", zap.Error(err))
				logger.Error("updated bootstrapper list is empty")
				logger.Error("updated bootstrapper list failed", zap.Error(err))
		logger.Fatal("could not create cluster topology watch", zap.Error(err))
		logger.Fatal("could not construct database", zap.Error(err))
		logger.Fatal("could not open database", zap.Error(err))
		logger.Fatal("could not open tchannelthrift interface",
			zap.String("address", cfg.ListenAddress), zap.Error(err))
	logger.Info("node tchannelthrift: listening", zap.String("address", cfg.ListenAddress))
		logger.Fatal("could not open tchannelthrift interface",
			zap.String("address", cfg.ClusterListenAddress), zap.Error(err))
	logger.Info("cluster tchannelthrift: listening", zap.String("address", cfg.ClusterListenAddress))
		logger.Fatal("could not open httpjson interface",
			zap.String("address", cfg.HTTPNodeListenAddress), zap.Error(err))
	logger.Info("node httpjson: listening", zap.String("address", cfg.HTTPNodeListenAddress))
		logger.Fatal("could not open httpjson interface",
			zap.String("address", cfg.HTTPClusterListenAddress), zap.Error(err))
	logger.Info("cluster httpjson: listening", zap.String("address", cfg.HTTPClusterListenAddress))
				logger.Error("debug server could not listen",
					zap.String("address", cfg.DebugListenAddress), zap.Error(err))
			logger.Fatal("could not bootstrap database", zap.Error(err))
		logger.Info("bootstrapped")
	logger.Warn("interrupt", zap.Error(interruptErr))
			logger.Error("close database error", zap.Error(err))
		logger.Info("server closed")
		logger.Error("server closed after timeout", zap.Duration("timeout", closeTimeout))
func bgValidateProcessLimits(logger *zap.Logger) {
		logger.Warn("cannot validate process limits: invalid configuration found",
			zap.String("message", message))
		logger.Warn("invalid configuration found, refer to linked documentation for more information",
			zap.String("url", xdocs.Path("operational_guide/kernel_configuration")),
			zap.Error(err),
		)
	logger *zap.Logger,
			logger.Warn("error resolving cluster new series insert limit", zap.Error(err))
		logger.Warn("unable to set cluster new series insert limit", zap.Error(err))
		logger.Error("could not watch cluster new series insert limit", zap.Error(err))
					logger.Warn("unable to parse new cluster new series insert limit", zap.Error(err))
				logger.Warn("unable to set cluster new series insert limit", zap.Error(err))
	logger *zap.Logger,
	logger *zap.Logger,
		logger.Error("could not resolve KV", zap.String("key", key), zap.Error(err))
			logger.Error("could not unmarshal KV key", zap.String("key", key), zap.Error(err))
			logger.Error("could not process value of KV", zap.String("key", key), zap.Error(err))
			logger.Info("set KV key", zap.String("key", key), zap.Any("value", protoValue.Value))
		logger.Error("could not watch KV key", zap.String("key", key), zap.Error(err))
					logger.Warn("could not set default for KV key", zap.String("key", key), zap.Error(err))
				logger.Warn("could not unmarshal KV key", zap.String("key", key), zap.Error(err))
				logger.Warn("could not process change for KV key", zap.String("key", key), zap.Error(err))
			logger.Info("set KV key", zap.String("key", key), zap.Any("value", protoValue.Value))
	logger *zap.Logger,
		logger.Fatal("could not watch value for key with KV",
			zap.String("key", kvconfig.BootstrapperKey))
				logger.Error("error converting KV update to string array",
					zap.String("key", kvconfig.BootstrapperKey),
					zap.Error(err),
				)
	logger *zap.Logger,
	schema *desc.MessageDescriptor,
		logger.Sugar().Infof("bytes pool registering bucket capacity=%d, size=%d, "+
			bucket.RefillLowWaterMarkOrDefault(), bucket.RefillHighWaterMarkOrDefault())
		logger.Fatal("unrecognized pooling type", zap.Any("type", policy.Type))
	logger.Sugar().Infof("bytes pool %s init", policy.Type)
		if schema != nil {
			enc := proto.NewEncoder(time.Time{}, encodingOpts)
			enc.SetSchema(schema)
			return enc
		}

		if schema != nil {
			return proto.NewIterator(r, schema, encodingOpts)
		}
	bucketPool := series.NewBufferBucketPool(
		poolOptions(policy.BufferBucketPool, scope.SubScope("buffer-bucket-pool")))
	bucketVersionsPool := series.NewBufferBucketVersionsPool(
		poolOptions(policy.BufferBucketVersionsPool, scope.SubScope("buffer-bucket-versions-pool")))

		SetWriteBatchPool(writeBatchPool).
		SetBufferBucketPool(bucketPool).
		SetBufferBucketVersionsPool(bucketVersionsPool)
		SetReaderIteratorPool(iteratorPool).
		SetMultiReaderIteratorPool(multiIteratorPool).