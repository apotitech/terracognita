package google

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cycloidio/terracognita/filter"
	"github.com/cycloidio/terracognita/provider"
)

// ResourceType is the type used to define all the Resources
// from the Provider
type ResourceType int

//go:generate enumer -type ResourceType -addprefix google_ -transform snake -linecomment
const (
	ComputeInstance ResourceType = iota
	ComputeFirewall
	ComputeNetwork
	// With Google, an HTTP(S) load balancer has 3 parts:
	// * backend configuration: instance_group, backend_service and health_check
	// * host and path rules: url_map
	// * frontend configuration: target_http(s)_proxy + global_forwarding_rule
	ComputeHealthCheck
	ComputeInstanceGroup
	ComputeInstanceIAMPolicy
	ComputeBackendBucket
	ComputeBackendService
	ComputeSSLCertificate
	ComputeTargetHTTPProxy
	ComputeTargetHTTPSProxy
	ComputeURLMap
	ComputeGlobalForwardingRule
	ComputeForwardingRule
	ComputeDisk
	DNSManagedZone
	DNSRecordSet
	ProjectIAMCustomRole
	StorageBucket
	StorageBucketIAMPolicy
	SQLDatabaseInstance

	noFilter = ""
)

type rtFn func(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error)

var (
	resources = map[ResourceType]rtFn{
		ComputeInstance:             computeInstance,
		ComputeFirewall:             computeFirewall,
		ComputeNetwork:              computeNetwork,
		ComputeHealthCheck:          computeHealthCheck,
		ComputeInstanceGroup:        computeInstanceGroup,
		ComputeInstanceIAMPolicy:    computeInstanceIAMPolicy,
		ComputeBackendService:       computeBackendService,
		ComputeBackendBucket:        computeBackendBucket,
		ComputeSSLCertificate:       computeSSLCertificate,
		ComputeTargetHTTPProxy:      computeTargetHTTPProxy,
		ComputeTargetHTTPSProxy:     computeTargetHTTPSProxy,
		ComputeURLMap:               computeURLMap,
		ComputeGlobalForwardingRule: computeGlobalForwardingRule,
		ComputeForwardingRule:       computeForwardingRule,
		ComputeDisk:                 computeDisk,
		DNSManagedZone:              managedZoneDNS,
		DNSRecordSet:                recordSetDNS,
		ProjectIAMCustomRole:        projectIAMCustomRole,
		StorageBucket:               storageBucket,
		StorageBucketIAMPolicy:      storageBucketIAMPolicy,
		SQLDatabaseInstance:         sqlDatabaseInstance,
	}
)

func initializeFilter(filters *filter.Filter) string {
	var b bytes.Buffer
	for _, t := range filters.Tags {
		// if multiple tags, we suppose it's a "AND" operation
		b.WriteString(fmt.Sprintf("(labels.%s=%s) ", t.Name, t.Value))
	}
	return b.String()
}

func computeInstance(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	f := initializeFilter(filters)
	instancesList, err := g.gcpr.ListInstances(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list instances from reader")
	}
	resources := make([]provider.Resource, 0)
	for z, instances := range instancesList {
		for _, instance := range instances {
			r := provider.NewResource(fmt.Sprintf("%s/%s/%s", g.Project(), z, instance.Name), resourceType, g)
			resources = append(resources, r)
		}
	}
	return resources, nil
}

func computeFirewall(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	firewalls, err := g.gcpr.ListFirewalls(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list firewalls from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, firewall := range firewalls {
		r := provider.NewResource(firewall.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeNetwork(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	networks, err := g.gcpr.ListNetworks(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list networks from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, network := range networks {
		r := provider.NewResource(network.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeHealthCheck(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	checks, err := g.gcpr.ListHealthChecks(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list health checks from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, check := range checks {
		r := provider.NewResource(check.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeInstanceGroup(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	instanceGroups, err := g.gcpr.ListInstanceGroups(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list instance groups from reader")
	}
	resources := make([]provider.Resource, 0)
	for z, groups := range instanceGroups {
		for _, group := range groups {
			r := provider.NewResource(fmt.Sprintf("%s/%s/%s", g.Project(), z, group.Name), resourceType, g)
			resources = append(resources, r)
		}
	}
	return resources, nil
}

func computeBackendService(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	backends, err := g.gcpr.ListBackendServices(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list backend services from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, backend := range backends {
		r := provider.NewResource(backend.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeURLMap(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	maps, err := g.gcpr.ListURLMaps(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list URL maps from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, urlMap := range maps {
		r := provider.NewResource(urlMap.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeTargetHTTPProxy(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	targets, err := g.gcpr.ListTargetHTTPProxies(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list target http proxies from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, target := range targets {
		r := provider.NewResource(target.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeTargetHTTPSProxy(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	targets, err := g.gcpr.ListTargetHTTPSProxies(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list target https proxies from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, target := range targets {
		r := provider.NewResource(target.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeSSLCertificate(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	certs, err := g.gcpr.ListSSLCertificates(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list SSL certificates from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, cert := range certs {
		r := provider.NewResource(cert.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeGlobalForwardingRule(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	f := initializeFilter(filters)
	rules, err := g.gcpr.ListGlobalForwardingRules(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list global forwarding rules from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, rule := range rules {
		r := provider.NewResource(rule.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeForwardingRule(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	f := initializeFilter(filters)
	rules, err := g.gcpr.ListForwardingRules(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list global forwarding rules from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, rule := range rules {
		r := provider.NewResource(rule.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func computeDisk(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	f := initializeFilter(filters)
	disksList, err := g.gcpr.ListDisks(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list disks from reader")
	}
	resources := make([]provider.Resource, 0)
	for z, disks := range disksList {
		for _, disk := range disks {
			r := provider.NewResource(fmt.Sprintf("%s/%s", z, disk.Name), resourceType, g)
			resources = append(resources, r)
		}
	}
	return resources, nil
}

func storageBucket(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	buckets, err := g.gcpr.ListBuckets(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list global forwarding rules from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, bucket := range buckets {
		r := provider.NewResource(bucket.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func sqlDatabaseInstance(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	instances, err := g.gcpr.ListStorageInstances(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list sql storage instances rules from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, instance := range instances {
		r := provider.NewResource(instance.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func managedZoneDNS(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	zones, err := g.gcpr.ListManagedZones(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list DNS managed zone from reader")
	}
	resources := make([]provider.Resource, 0)
	for _, zone := range zones {
		r := provider.NewResource(zone.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func recordSetDNS(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	managedZones, err := managedZoneDNS(ctx, g, resourceType, filters)
	if err != nil {
		return nil, errors.Wrap(err, "unable to previously fetch managed zones")
	}
	zones := make([]string, 0, len(managedZones))
	for _, zone := range managedZones {
		zones = append(zones, zone.ID())
	}
	rrsetsList, err := g.gcpr.ListResourceRecordSets(ctx, zones)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list resources record se record sett from reader")
	}
	resources := make([]provider.Resource, 0)
	for z, rrsets := range rrsetsList {
		for _, rrset := range rrsets {
			r := provider.NewResource(fmt.Sprintf("%s/%s/%s", z, rrset.Name, rrset.Type), resourceType, g)
			resources = append(resources, r)
		}
	}
	return resources, nil
}

func computeBackendBucket(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	backends, err := g.gcpr.ListBackendBuckets(ctx, noFilter)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list backend buckets from reader")
	}
	resources := make([]provider.Resource, 0, len(backends))
	for _, backend := range backends {
		r := provider.NewResource(backend.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

func projectIAMCustomRole(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	roles, err := g.gcpr.ListProjectIAMCustomRoles(ctx, fmt.Sprintf("projects/%s", g.gcpr.project))
	if err != nil {
		return nil, errors.Wrap(err, "unable to list project IAM custom roles from reader")
	}
	resources := make([]provider.Resource, 0, len(roles))
	for _, role := range roles {
		r := provider.NewResource(role.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

// storageBucketIAMPolicy will import the policies binded to a bucket. We need to iterate over the
// bucket list
func storageBucketIAMPolicy(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	buckets, err := g.gcpr.ListBuckets(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list bucket policies custom roles from reader")
	}
	resources := make([]provider.Resource, 0, len(buckets))
	for _, bucket := range buckets {
		r := provider.NewResource(bucket.Name, resourceType, g)
		resources = append(resources, r)
	}
	return resources, nil
}

// computeInstanceIAMPolicy will import the policies binded to a compute instance. We need to iterate over the
// compute instance list
func computeInstanceIAMPolicy(ctx context.Context, g *google, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	f := initializeFilter(filters)
	list, err := g.gcpr.ListInstances(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list compute instances from reader")
	}
	resources := make([]provider.Resource, 0)
	for zone, instances := range list {
		for _, instance := range instances {
			r := provider.NewResource(fmt.Sprintf("projects/%s/zones/%s/instances/%s", g.Project(), zone, instance.Name), resourceType, g)
			resources = append(resources, r)
		}
	}
	return resources, nil
}
