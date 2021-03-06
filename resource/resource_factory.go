package resource

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/worker"
)

//go:generate counterfeiter . ResourceFactoryFactory

type ResourceFactoryFactory interface {
	FactoryFor(workerClient worker.Client) ResourceFactory
}

type resourceFactoryFactory struct{}

func NewResourceFactoryFactory() ResourceFactoryFactory {
	return &resourceFactoryFactory{}
}

func (f *resourceFactoryFactory) FactoryFor(workerClient worker.Client) ResourceFactory {
	return &resourceFactory{
		workerClient: workerClient,
	}
}

//go:generate counterfeiter . ResourceFactory

type ResourceFactory interface {
	NewResource(
		logger lager.Logger,
		id worker.Identifier,
		metadata worker.Metadata,
		resourceSpec worker.ContainerSpec,
		resourceTypes atc.ResourceTypes,
		imageFetchingDelegate worker.ImageFetchingDelegate,
		resourceSources map[string]worker.ArtifactSource,
	) (Resource, []string, error)

	NewBuildResource(
		logger lager.Logger,
		id worker.Identifier,
		metadata worker.Metadata,
		containerSpec worker.ContainerSpec,
		resourceTypes atc.ResourceTypes,
		imageFetchingDelegate worker.ImageFetchingDelegate,
		inputSources []InputSource,
		outputPaths map[string]string,
	) (Resource, []InputSource, error)

	NewCheckResource(
		logger lager.Logger,
		id worker.Identifier,
		metadata worker.Metadata,
		resourceSpec worker.ContainerSpec,
		resourceTypes atc.ResourceTypes,
	) (Resource, error)

	NewCheckResourceForResourceType(
		logger lager.Logger,
		id worker.Identifier,
		metadata worker.Metadata,
		resourceSpec worker.ContainerSpec,
		resourceTypes atc.ResourceTypes,
	) (Resource, error)
}

//go:generate counterfeiter . InputSource

type InputSource interface {
	Name() worker.ArtifactName
	Source() worker.ArtifactSource
	MountPath() string
}

type resourceFactory struct {
	workerClient worker.Client
}

func (f *resourceFactory) NewResource(
	logger lager.Logger,
	id worker.Identifier,
	metadata worker.Metadata,
	resourceSpec worker.ContainerSpec,
	resourceTypes atc.ResourceTypes,
	imageFetchingDelegate worker.ImageFetchingDelegate,
	resourceSources map[string]worker.ArtifactSource,
) (Resource, []string, error) {
	container, missingSourceNames, err := f.workerClient.FindOrCreateContainerForIdentifier(
		logger,
		id,
		metadata,
		resourceSpec,
		resourceTypes,
		imageFetchingDelegate,
		resourceSources,
	)
	if err != nil {
		return nil, nil, err
	}

	return NewResourceForContainer(container), missingSourceNames, nil
}

func (f *resourceFactory) NewBuildResource(
	logger lager.Logger,
	id worker.Identifier,
	metadata worker.Metadata,
	containerSpec worker.ContainerSpec,
	resourceTypes atc.ResourceTypes,
	imageFetchingDelegate worker.ImageFetchingDelegate,
	inputSources []InputSource,
	outputPaths map[string]string,
) (Resource, []InputSource, error) {
	compatibleWorkers, err := f.workerClient.AllSatisfying(containerSpec.WorkerSpec(), resourceTypes)
	if err != nil {
		return nil, nil, err
	}

	// find the worker with the most volumes
	mounts := []worker.VolumeMount{}
	missingSources := []InputSource{}
	var chosenWorker worker.Worker

	for _, w := range compatibleWorkers {
		candidateMounts := []worker.VolumeMount{}
		missing := []InputSource{}

		for _, inputSource := range inputSources {
			ourVolume, found, err := inputSource.Source().VolumeOn(w)
			if err != nil {
				return nil, nil, err
			}

			if found {
				candidateMounts = append(candidateMounts, worker.VolumeMount{
					Volume:    ourVolume,
					MountPath: inputSource.MountPath(),
				})
			} else {
				missing = append(missing, inputSource)
			}
		}

		if len(candidateMounts) >= len(mounts) {
			mounts = candidateMounts
			missingSources = missing
			chosenWorker = w
		}
	}

	containerSpec.Inputs = mounts

	container, err := chosenWorker.FindOrCreateBuildContainer(
		logger,
		nil,
		imageFetchingDelegate,
		id,
		metadata,
		containerSpec,
		resourceTypes,
		outputPaths,
	)
	if err != nil {
		return nil, nil, err
	}

	return NewResourceForContainer(container), missingSources, nil
}

func (f *resourceFactory) NewCheckResource(
	logger lager.Logger,
	id worker.Identifier,
	metadata worker.Metadata,
	resourceSpec worker.ContainerSpec,
	resourceTypes atc.ResourceTypes,
) (Resource, error) {
	container, err := f.workerClient.FindOrCreateResourceCheckContainer(
		logger,
		nil,
		worker.NoopImageFetchingDelegate{},
		id,
		metadata,
		resourceSpec,
		resourceTypes,
		id.CheckType,
		id.CheckSource,
	)
	if err != nil {
		return nil, err
	}

	return NewResourceForContainer(container), nil
}

func (f *resourceFactory) NewCheckResourceForResourceType(
	logger lager.Logger,
	id worker.Identifier,
	metadata worker.Metadata,
	resourceSpec worker.ContainerSpec,
	resourceTypes atc.ResourceTypes,
) (Resource, error) {
	container, err := f.workerClient.FindOrCreateResourceTypeCheckContainer(
		logger,
		nil,
		worker.NoopImageFetchingDelegate{},
		id,
		metadata,
		resourceSpec,
		resourceTypes,
		id.CheckType,
		id.CheckSource,
	)
	if err != nil {
		return nil, err
	}

	return NewResourceForContainer(container), nil
}
