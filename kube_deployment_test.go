package test
import (
  "testing"
  "path/filepath"
  "github.com/stretchr/testify/require"
  "github.com/gruntwork-io/terratest/modules/k8s"
  "time"
  "fmt"
  http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)
func TestKubeDeployment(t *testing.T) {
  t.Parallel()
  options := k8s.NewKubectlOptions("", "", "default")
  pathDeployment, err := filepath.Abs("deploy.yaml")
  require.NoError(t, err)
  pathService, err := filepath.Abs("svc.yaml")
  require.NoError(t, err)
  defer k8s.KubectlDelete(t, options, pathDeployment)
  defer k8s.KubectlDelete(t, options, pathService)
  k8s.KubectlApply(t, options, pathDeployment)
  k8s.KubectlApply(t, options, pathService)
  service := k8s.GetService(t, options, "nginx-svc")
  require.Equal(t, service.Name, "nginx-svc")
  k8s.WaitUntilServiceAvailable(t, options, "nginx-svc", 10, 1 * time.Second)
  //service := k8s.GetService(t, options, "nginx-svc")
  url := fmt.Sprintf("http://%s", k8s.GetServiceEndpoint(t, options, service, 80))
  http_helper.HttpGetWithRetry(t, url, nil, 200, "<html>\n<body>\n<p>Helloworld !</p>\n</body>\n</html>", 30, 3 * time.Second)
}
