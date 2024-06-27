package podmonitor

import (
	//"context"

	synccontext "github.com/loft-sh/vcluster/pkg/controllers/syncer/context"
	"github.com/loft-sh/vcluster/pkg/controllers/syncer/translator"
	"github.com/loft-sh/vcluster/pkg/scheme"
	synctypes "github.com/loft-sh/vcluster/pkg/types"
	"github.com/loft-sh/vcluster/pkg/util/translate"
	"k8s.io/apimachinery/pkg/api/equality"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	internaltranslators "github.com/codefresh-contrib/vcluster-prometheus-operator-plugin/v2/translators"
)

func init() {
	// Make sure our scheme is registered
	_ = promoperatorv1.AddToScheme(scheme.Scheme)
}

func NewPodMonitorSyncer(ctx *synccontext.RegisterContext) synctypes.Base {
	return &podMonitorSyncer{
		NamespacedTranslator: translator.NewNamespacedTranslator(ctx, "podmonitor", &promoperatorv1.PodMonitor{}),
	}
}

type podMonitorSyncer struct {
	translator.NamespacedTranslator
}

var _ synctypes.Initializer = &podMonitorSyncer{}

func (s *podMonitorSyncer) Init(ctx *synccontext.RegisterContext) error {
	_, _, err := translate.EnsureCRDFromPhysicalCluster(ctx.Context, ctx.PhysicalManager.GetConfig(), ctx.VirtualManager.GetConfig(), promoperatorv1.SchemeGroupVersion.WithKind("PodMonitor"))
	return err
}

var _ synctypes.Syncer = &podMonitorSyncer{}

func (s *podMonitorSyncer) SyncToHost(ctx *synccontext.SyncContext, vObj client.Object) (ctrl.Result, error) {
	return s.SyncToHostCreate(ctx, vObj, s.TranslateMetadata(ctx.Context, vObj).(*promoperatorv1.PodMonitor))
}

func (s *podMonitorSyncer) Sync(ctx *synccontext.SyncContext, pObj client.Object, vObj client.Object) (ctrl.Result, error) {
	return s.SyncToHostUpdate(ctx, vObj, s.translateUpdate(ctx, pObj.(*promoperatorv1.PodMonitor), vObj.(*promoperatorv1.PodMonitor)))
}

func (s *podMonitorSyncer) translateUpdate(ctx *synccontext.SyncContext, pObj, vObj *promoperatorv1.PodMonitor) *promoperatorv1.PodMonitor {
	var updated *promoperatorv1.PodMonitor

	// check annotations & labels
	changed, updatedAnnotations, updatedLabels := s.TranslateMetadataUpdate(ctx.Context, vObj, pObj)

	if changed {
		updated = translator.NewIfNil(updated, pObj)
		updated.Labels = updatedLabels
		updated.Annotations = updatedAnnotations
	}

	newSpec := internaltranslators.TranslatePodMonitorSpec(&vObj.Spec, vObj.Namespace)

	// check spec
	if !equality.Semantic.DeepEqual(newSpec, pObj.Spec) {
		updated = translator.NewIfNil(updated, pObj)
		updated.Spec = *newSpec
	}

	return updated
}
