package persist

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	sharedutil "github.com/nholuongut/argo-events/pkg/shared/util"
)

type EventPersist interface {
	Save(event *Event) error
	Get(key string) (*Event, error)
	IsEnabled() bool
}

type Event struct {
	EventKey     string
	EventPayload string
}

type ConfigMapPersist struct {
	ctx              context.Context
	kubeClient       kubernetes.Interface
	name             string
	namespace        string
	createIfNotExist bool
}

func createConfigmap(ctx context.Context, client kubernetes.Interface, name, namespace string) (*v1.ConfigMap, error) {
	cm := v1.ConfigMap{}
	cm.Name = name
	cm.Namespace = namespace
	return client.CoreV1().ConfigMaps(namespace).Create(ctx, &cm, metav1.CreateOptions{})
}

func NewConfigMapPersist(ctx context.Context, client kubernetes.Interface, configmap *v1alpha1.ConfigMapPersistence, namespace string) (EventPersist, error) {
	if configmap == nil {
		return nil, fmt.Errorf("persistence configuration is nil")
	}
	_, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, configmap.Name, metav1.GetOptions{})
	if err != nil {
		if apierr.IsNotFound(err) && configmap.CreateIfNotExist {
			_, err = createConfigmap(ctx, client, configmap.Name, namespace)
			if err != nil {
				if !apierr.IsAlreadyExists(err) {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}
	cmp := ConfigMapPersist{
		ctx:              ctx,
		kubeClient:       client,
		name:             configmap.Name,
		namespace:        namespace,
		createIfNotExist: configmap.CreateIfNotExist,
	}
	return &cmp, nil
}

func (cmp *ConfigMapPersist) IsEnabled() bool {
	return true
}

func (cmp *ConfigMapPersist) Save(event *Event) error {
	if event == nil {
		return fmt.Errorf("event object is nil")
	}
	// Using Connect util func for backoff retry if K8s API returns error
	err := sharedutil.DoWithRetry(&sharedutil.DefaultBackoff, func() error {
		cm, err := cmp.kubeClient.CoreV1().ConfigMaps(cmp.namespace).Get(cmp.ctx, cmp.name, metav1.GetOptions{})
		if err != nil {
			if apierr.IsNotFound(err) && cmp.createIfNotExist {
				cm, err = createConfigmap(cmp.ctx, cmp.kubeClient, cmp.name, cmp.namespace)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		if len(cm.Data) == 0 {
			cm.Data = make(map[string]string)
		}

		cm.Data[event.EventKey] = event.EventPayload
		_, err = cmp.kubeClient.CoreV1().ConfigMaps(cmp.namespace).Update(cmp.ctx, cm, metav1.UpdateOptions{})

		return err
	})

	if err != nil {
		return err
	}
	return nil
}

func (cmp *ConfigMapPersist) Get(key string) (*Event, error) {
	cm, err := cmp.kubeClient.CoreV1().ConfigMaps(cmp.namespace).Get(cmp.ctx, cmp.name, metav1.GetOptions{})
	if err != nil {
		if apierr.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	payload, exist := cm.Data[key]
	if !exist {
		return nil, nil
	}
	return &Event{EventKey: key, EventPayload: payload}, nil
}

type NullPersistence struct {
}

func (n *NullPersistence) Save(event *Event) error {
	return nil
}

func (n *NullPersistence) Get(key string) (*Event, error) {
	return nil, nil
}

func (cmp *NullPersistence) IsEnabled() bool {
	return false
}
