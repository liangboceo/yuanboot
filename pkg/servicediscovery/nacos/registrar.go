package nacos

import (
	"errors"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/servicediscovery"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	sd "github.com/liangboceo/yuanboot/pkg/servicediscovery"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"strings"
	"time"
)

type Registrar struct {
	cacheLocalInstance servicediscovery.ServiceInstance
	logger             xlog.ILogger
	config             *Config
	client             naming_client.INamingClient
}

func NewServerDiscoveryWithDI(configuration abstractions.IConfiguration, env *abstractions.HostEnvironment) servicediscovery.IServiceDiscovery {
	sdType, ok := configuration.Get("yuanboot.cloud.discovery.type").(string)
	if !ok || sdType != "nacos" {
		panic(errors.New("yuanboot.cloud.discovery.type is not config node"))
	}
	path, _ := configuration.Get("yuanboot.application.server.path").(string)
	section := configuration.GetSection("yuanboot.cloud.discovery.metadata")
	if section == nil {
		panic(errors.New("yuanboot.cloud.discovery.metadata is not config node"))
	}
	option := &Config{}
	section.Unmarshal(&option)
	if option.GroupName == "" {
		option.GroupName = GroupName
	}
	if option.Cluster == "" {
		option.Cluster = Cluster
	}
	if option.Path == "" {
		option.Path = path
	}
	option.ENV = env

	return NewServerDiscovery(option)
}

func NewServerDiscovery(option *Config) servicediscovery.IServiceDiscovery {
	logger := xlog.GetXLogger("Server Discovery nacos")
	nacosRegister := &Registrar{}
	var serverConfigs []constant.ServerConfig
	urls := strings.Split(option.Url, ";")
	for _, url := range urls {
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			ContextPath: "/nacos",
			IpAddr:      url,
			Port:        option.Port,
		})
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         option.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "info",
		CacheDir:            logger.GetLogPath(),
		CustomLogger:        logger,
	}

	if option.Auth != nil && option.Auth.Enable {
		clientConfig.Username = option.Auth.User
		clientConfig.Password = option.Auth.Password
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	nacosRegister.client = namingClient
	nacosRegister.config = option
	nacosRegister.logger = logger

	logger.Debugf("url:%s, namespace:%s , group:%s , cluster:%s ;", option.Url, option.NamespaceId, option.GroupName, option.Cluster)
	return nacosRegister
}

func (registrar Registrar) GetName() string {
	return "nacos"
}

func (registrar *Registrar) Register() error {
	registrar.cacheLocalInstance = sd.CreateServiceInstance(registrar.config.ENV)
	success, err := registrar.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          registrar.cacheLocalInstance.GetHost(),
		Port:        registrar.cacheLocalInstance.GetPort(),
		ServiceName: registrar.cacheLocalInstance.GetServiceName(),
		Weight:      10,
		ClusterName: registrar.config.Cluster,
		GroupName:   registrar.config.GroupName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"yuanboot_version":                  registrar.config.ENV.Version,
			"yuanboot_env":                      registrar.config.ENV.Profile,
			"yuanboot_author":                   "jackson",
			"yuanboot_application_description":  "quickly,quickly,quickly,quickly,quickly!!!!!",
			"yuanboot_application_context-path": registrar.config.Path,
			"yuanboot_application_pid":          strconv.Itoa(registrar.config.ENV.PID),
			"yuanboot_application_create_time":  time.Now().Format("2006-01-02 15:04:05"),
			"yuanboot_application_name":         registrar.config.ENV.ApplicationName,
		},
	})
	if err != nil {
		registrar.logger.Error(err.Error())
	}
	registrar.logger.Debugf("Registrar IP: %s , Success: %v", registrar.config.ENV.Host, success)
	return err
}

func (registrar *Registrar) Update() error {

	return nil
}

func (registrar *Registrar) Unregister() error {
	if registrar.cacheLocalInstance == nil {
		return nil
	}
	_, err := registrar.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          registrar.cacheLocalInstance.GetHost(),
		Port:        registrar.cacheLocalInstance.GetPort(),
		Cluster:     registrar.cacheLocalInstance.GetClusterName(),
		ServiceName: registrar.cacheLocalInstance.GetServiceName(),
		GroupName:   registrar.cacheLocalInstance.GetGroupName(),
		Ephemeral:   true,
	})
	if err != nil {
		registrar.logger.Error(err.Error())
	}
	return err
}

func (registrar *Registrar) GetHealthyInstances(serviceName string) []servicediscovery.ServiceInstance {
	// SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	instances, err := registrar.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   registrar.config.GroupName,         // default value is DEFAULT_GROUP
		Clusters:    []string{registrar.config.Cluster}, // default value is DEFAULT
		HealthyOnly: true,
	})
	if err != nil {
		return nil
	}
	return convInstance(registrar.config.GroupName, instances)
}

func (registrar *Registrar) GetAllInstances(serviceName string) []servicediscovery.ServiceInstance {
	instances, err := registrar.client.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   registrar.config.GroupName,         // default value is DEFAULT_GROUP
		Clusters:    []string{registrar.config.Cluster}, // default value is DEFAULT
	})

	if err != nil {
		return nil
	}
	return convInstance(registrar.config.GroupName, instances)
}

func convInstance(groupName string, sourceInstances []model.Instance) []servicediscovery.ServiceInstance {
	var serviceList []servicediscovery.ServiceInstance
	for _, s := range sourceInstances {
		instance := &servicediscovery.DefaultServiceInstance{
			Id:          s.InstanceId,
			ServiceName: s.ServiceName,
			Host:        s.Ip,
			Port:        s.Port,
			ClusterName: s.ClusterName,
			GroupName:   groupName,
			Enable:      true,
			Weight:      s.Weight,
			Healthy:     s.Healthy,
			Metadata:    s.Metadata,
		}
		serviceList = append(serviceList, instance)
	}
	return serviceList
}

func (registrar *Registrar) Destroy() error {
	registrar.logger.Debugf("Destroy")
	err := registrar.Unregister()
	return err
}

func (registrar *Registrar) Watch(opts ...servicediscovery.WatchOption) (servicediscovery.Watcher, error) {
	return newWatcher(registrar.client, registrar.config, registrar.logger, opts...)
}

func (registrar *Registrar) GetAllServices() ([]*servicediscovery.Service, error) {
	serviceList, _ := registrar.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: registrar.config.NamespaceId,
		GroupName: registrar.config.GroupName,
		PageNo:    1,
		PageSize:  1000,
	})

	services := make([]*servicediscovery.Service, 0)
	for _, serviceName := range serviceList.Doms {
		services = append(services, &servicediscovery.Service{Name: serviceName})
	}
	return services, nil
}
