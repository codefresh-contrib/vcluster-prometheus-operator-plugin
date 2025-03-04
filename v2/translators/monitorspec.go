package translators

import (
	"github.com/loft-sh/vcluster/pkg/util/translate"
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MonitorSpec struct {
	NamespaceSelector promoperatorv1.NamespaceSelector
	Selector          metav1.LabelSelector
	JobLabel          string
	Namespace         string
}

func (s *MonitorSpec) rewriteSpec() *MonitorSpec {

	newSpec := s
	monitorNamespaces := []string{}

	// If not any namespace - it's required to constuct a list of namespaces
	if !s.NamespaceSelector.Any {
		if len(s.NamespaceSelector.MatchNames) > 0 {
			monitorNamespaces = s.NamespaceSelector.DeepCopy().MatchNames
		} else {
			monitorNamespaces = append(monitorNamespaces, s.Namespace)
		}
	}

	// Clear namespace selector as it does not apply on host cluster
	newSpec.NamespaceSelector = promoperatorv1.NamespaceSelector{}
	newSpec.Selector = *translate.Default.TranslateLabelSelector(&s.Selector)

	if len(monitorNamespaces) > 0 {
		nsExpression := metav1.LabelSelectorRequirement{Key: translate.NamespaceLabel, Operator: metav1.LabelSelectorOpIn, Values: monitorNamespaces}
		newSpec.Selector.MatchExpressions = append(s.Selector.MatchExpressions, nsExpression)
	}

	// Translate job labels
	if len(s.JobLabel) > 0 {
		newSpec.JobLabel = translate.Default.ConvertLabelKey(s.JobLabel)
	}

	return newSpec
}

func TranslatePodMonitorSpec(vPodMonitorSpec *promoperatorv1.PodMonitorSpec, podMonitorNamespace string) *promoperatorv1.PodMonitorSpec {
	changedSpec := &MonitorSpec{Selector: vPodMonitorSpec.Selector, NamespaceSelector: vPodMonitorSpec.NamespaceSelector, JobLabel: vPodMonitorSpec.JobLabel, Namespace: podMonitorNamespace}
	changedSpec = changedSpec.rewriteSpec()
	newSpec := vPodMonitorSpec.DeepCopy()

	newSpec.Selector = changedSpec.Selector
	newSpec.NamespaceSelector = changedSpec.NamespaceSelector
	newSpec.JobLabel = changedSpec.JobLabel

	return newSpec
}

func TranslateServiceMonitorSpec(vPodMonitorSpec *promoperatorv1.ServiceMonitorSpec, serviceMonitorNamespace string) *promoperatorv1.ServiceMonitorSpec {
	changedSpec := &MonitorSpec{Selector: vPodMonitorSpec.Selector, NamespaceSelector: vPodMonitorSpec.NamespaceSelector, JobLabel: vPodMonitorSpec.JobLabel, Namespace: serviceMonitorNamespace}
	changedSpec = changedSpec.rewriteSpec()
	newSpec := vPodMonitorSpec.DeepCopy()

	newSpec.Selector = changedSpec.Selector
	newSpec.NamespaceSelector = changedSpec.NamespaceSelector
	newSpec.JobLabel = changedSpec.JobLabel

	return newSpec
}
