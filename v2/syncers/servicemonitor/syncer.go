package servicemonitor

import (
	_ "embed"
	"fmt"

	"github.com/loft-sh/vcluster/pkg/mappings/generic"
	"github.com/loft-sh/vcluster/pkg/patcher"
	"github.com/loft-sh/vcluster/pkg/scheme"
	"github.com/loft-sh/vcluster/pkg/syncer"
	"github.com/loft-sh/vcluster/pkg/syncer/synccontext"
	"github.com/loft-sh/vcluster/pkg/syncer/translator"
	syncertypes "github.com/loft-sh/vcluster/pkg/syncer/types"
	"github.com/loft-sh/vcluster/pkg/util"
	"github.com/loft-sh/vcluster/pkg/util/translate"
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	// Make sure our scheme is registered
	_ = promoperatorv1.AddToScheme(scheme.Scheme)
}

// TODO: rewrite using v2/vendor/k8s.io/apimachinery/pkg/runtime/serializer/json/json.go
//
//go:embed crd-servicemonitors.yaml
var serviceMonitorCRD string

func NewServiceMonitorSyncer(ctx *synccontext.RegisterContext) (syncertypes.Base, error) {
	CRDManifest := []byte(serviceMonitorCRD)
	GVK := promoperatorv1.SchemeGroupVersion.WithKind("ServiceMonitor")

	err := util.EnsureCRD(ctx.Context, ctx.PhysicalManager.GetConfig(), CRDManifest, GVK)
	if err != nil {
		return nil, err
	}

	err = util.EnsureCRD(ctx.Context, ctx.VirtualManager.GetConfig(), CRDManifest, GVK)
	if err != nil {
		return nil, err
	}

	mapper, err := generic.NewMapper(ctx, &promoperatorv1.ServiceMonitor{}, translate.Default.HostName)
	if err != nil {
		return nil, err
	}

	return &serviceMonitorSyncer{
		GenericTranslator: translator.NewGenericTranslator(ctx, "servicemonitor", &promoperatorv1.ServiceMonitor{}, mapper),
	}, nil
}

type serviceMonitorSyncer struct {
	syncertypes.GenericTranslator
}

var _ syncertypes.Syncer = &serviceMonitorSyncer{}

func (s *serviceMonitorSyncer) Syncer() syncertypes.Sync[client.Object] {
	return syncer.ToGenericSyncer(s)
}

func (s *serviceMonitorSyncer) SyncToHost(ctx *synccontext.SyncContext, event *synccontext.SyncToHostEvent[*promoperatorv1.ServiceMonitor]) (ctrl.Result, error) {
	pObj := translate.HostMetadata(event.Virtual, s.VirtualToHost(ctx, types.NamespacedName{Name: event.Virtual.Name, Namespace: event.Virtual.Namespace}, event.Virtual))
	return patcher.CreateHostObject(ctx, event.Virtual, pObj, s.EventRecorder(), true)
}

func (s *serviceMonitorSyncer) Sync(ctx *synccontext.SyncContext, event *synccontext.SyncEvent[*promoperatorv1.ServiceMonitor]) (_ ctrl.Result, retErr error) {
	patchHelper, err := patcher.NewSyncerPatcher(ctx, event.Host, event.Virtual)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("new syncer patcher: %w", err)
	}

	defer func() {
		if err := patchHelper.Patch(ctx, event.Host, event.Virtual); err != nil {
			retErr = errors.NewAggregate([]error{retErr, err})
		}
		if retErr != nil {
			s.EventRecorder().Eventf(event.Virtual, "Warning", "SyncError", "Error syncing: %v", retErr)
		}
	}()

	// any changes made below here are correctly synced

	// sync metadata
	event.Host.Annotations = translate.HostAnnotations(event.Virtual, event.Host)
	event.Host.Labels = translate.HostLabels(event.Virtual, event.Host)

	// sync virtual to host
	event.Host.Spec = event.Virtual.Spec

	return ctrl.Result{}, nil
}

func (s *serviceMonitorSyncer) SyncToVirtual(ctx *synccontext.SyncContext, event *synccontext.SyncToVirtualEvent[*promoperatorv1.ServiceMonitor]) (ctrl.Result, error) {
	// virtual object is not here anymore, so we delete
	return patcher.DeleteHostObject(ctx, event.Host, nil, "virtual object was deleted")
}
