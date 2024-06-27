package podmonitor

import (
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/loft-sh/vcluster-sdk/plugin"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	"github.com/loft-sh/vcluster-sdk/translate"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/apimachinery/pkg/api/equality"
	internaltranslators "github.com/codefresh-contrib/vcluster-prometheus-operator-plugin/v1/translators"
)

func init() {
	// Make sure our scheme is registered
	_ = promoperatorv1.AddToScheme(plugin.Scheme)
}

func NewPodMonitorSyncer(ctx *synccontext.RegisterContext) syncer.Base {
	return &podMonitorSyncer{
		NamespacedTranslator: translator.NewNamespacedTranslator(ctx, "podmonitor", &promoperatorv1.PodMonitor{}),
	}
}

type podMonitorSyncer struct {
	translator.NamespacedTranslator
}

var _ syncer.Initializer = &podMonitorSyncer{}

func (s *podMonitorSyncer) Init(ctx *synccontext.RegisterContext) error {
	return translate.EnsureCRDFromPhysicalCluster(ctx.Context, ctx.PhysicalManager.GetConfig(), ctx.VirtualManager.GetConfig(), promoperatorv1.SchemeGroupVersion.WithKind("PodMonitor"))
}

func (s *podMonitorSyncer) SyncDown(ctx *synccontext.SyncContext, vObj client.Object) (ctrl.Result, error) {
	return s.SyncDownCreate(ctx, vObj, s.TranslateMetadata(vObj).(*promoperatorv1.PodMonitor))
}

func (s *podMonitorSyncer) Sync(ctx *synccontext.SyncContext, pObj client.Object, vObj client.Object) (ctrl.Result, error) {
	return s.SyncDownUpdate(ctx, vObj, s.translateUpdate(pObj.(*promoperatorv1.PodMonitor), vObj.(*promoperatorv1.PodMonitor)))
}

func newIfNil(updated *promoperatorv1.PodMonitor, pObj *promoperatorv1.PodMonitor) *promoperatorv1.PodMonitor {
	if updated == nil {
		return pObj.DeepCopy()
	}
	return updated
}

func (s *podMonitorSyncer) translateUpdate(pObj, vObj *promoperatorv1.PodMonitor) *promoperatorv1.PodMonitor {
	var updated *promoperatorv1.PodMonitor

	// check annotations & labels
	changed, updatedAnnotations, updatedLabels := s.TranslateMetadataUpdate(vObj, pObj)

	if changed {
		updated = newIfNil(updated, pObj)
		updated.Labels = updatedLabels
		updated.Annotations = updatedAnnotations
	}

	newSpec := internaltranslators.TranslatePodMonitorSpec(vObj.Spec.DeepCopy(), vObj.Namespace)

	if !equality.Semantic.DeepEqual(newSpec, pObj.Spec) {
		updated = newIfNil(updated, pObj)
		updated.Spec = *newSpec
	}

	return updated
}
